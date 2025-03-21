package sqlsubscriptionstore

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.skia.org/infra/go/sql/pool"
	"go.skia.org/infra/perf/go/sql/sqltest"
	"go.skia.org/infra/perf/go/subscription"
	pb "go.skia.org/infra/perf/go/subscription/proto/v1"
)

func setUp(t *testing.T) (subscription.Store, pool.Pool) {
	db := sqltest.NewSpannerDBForTests(t, "subscriptionstore")
	store, err := New(db)
	require.NoError(t, err)

	return store, db
}

// Insert two valid subscriptions and ensure they're queryable.
func TestInsert_ValidSubscriptions(t *testing.T) {
	ctx := context.Background()
	store, db := setUp(t)

	s := []*pb.Subscription{
		{
			Name:         "Test Subscription 1",
			Revision:     "abcd",
			BugLabels:    []string{"A", "B"},
			Hotlists:     []string{"C", "D"},
			BugComponent: "Component1>Subcomponent1",
			BugPriority:  1,
			BugSeverity:  2,
			BugCcEmails: []string{
				"abcd@efg.com",
				"1234@567.com",
			},
			ContactEmail: "test@owner.com",
		},
		{
			Name:         "Test Subscription 2",
			Revision:     "abcd",
			BugLabels:    []string{"1", "2"},
			Hotlists:     []string{"3", "4"},
			BugComponent: "Component2>Subcomponent2",
			BugPriority:  1,
			BugSeverity:  2,
			BugCcEmails: []string{
				"abcd@efg.com",
				"1234@567.com",
			},
			ContactEmail: "test@owner.com",
		},
	}
	tx, err := db.Begin(ctx)
	require.NoError(t, err)

	err = store.InsertSubscriptions(ctx, s, tx)
	require.NoError(t, err)

	err = tx.Commit(ctx)
	require.NoError(t, err)

	actual := getSubscriptionsFromDb(t, ctx, db)
	assert.ElementsMatch(t, actual, s)
}

// Test inserting two subscriptions with same primary key. The transaction
// should fail on second insert and no subscriptions should live in the DB
// due to transaction rollback.
func TestInsert_DuplicateSubscriptionKeys(t *testing.T) {
	ctx := context.Background()
	store, db := setUp(t)

	s := []*pb.Subscription{
		{
			Name:         "Test Subscription 1",
			Revision:     "abcd",
			BugLabels:    []string{"A", "B"},
			Hotlists:     []string{"C", "D"},
			BugComponent: "Component1>Subcomponent1",
			BugPriority:  1,
			BugSeverity:  2,
			BugCcEmails: []string{
				"abcd@efg.com",
				"1234@567.com",
			},
			ContactEmail: "test@owner.com",
		},
		{
			Name:         "Test Subscription 1",
			Revision:     "abcd",
			BugLabels:    []string{"1", "2"},
			Hotlists:     []string{"3", "4"},
			BugComponent: "Component2>Subcomponent2",
			BugPriority:  1,
			BugSeverity:  2,
			BugCcEmails: []string{
				"abcd@efg.com",
				"1234@567.com",
			},
			ContactEmail: "test@owner.com",
		},
	}

	tx, err := db.Begin(ctx)
	require.NoError(t, err)

	err = store.InsertSubscriptions(ctx, s, tx)
	require.Error(t, err)

	err = tx.Rollback(ctx)
	require.NoError(t, err)

	actual := getSubscriptionsFromDb(t, ctx, db)
	assert.Empty(t, actual)
}

// Test inserting two subscriptions with the same name but different revisions.
// The transaction should succeed.
func TestInsert_SameNameDifferentRevision(t *testing.T) {
	ctx := context.Background()
	store, db := setUp(t)

	s1 := &pb.Subscription{
		Name:         "Test Subscription",
		Revision:     "abcd",
		BugLabels:    []string{"A", "B"},
		Hotlists:     []string{"C", "D"},
		BugComponent: "Component1>Subcomponent1",
		BugPriority:  1,
		BugSeverity:  2,
		BugCcEmails: []string{
			"abcd@efg.com",
			"1234@567.com",
		},
		ContactEmail: "test@owner.com",
	}
	s2 := &pb.Subscription{
		Name:         "Test Subscription",
		Revision:     "efgh",
		BugLabels:    []string{"1", "2"},
		Hotlists:     []string{"3", "4"},
		BugComponent: "Component2>Subcomponent2",
		BugPriority:  1,
		BugSeverity:  2,
		BugCcEmails: []string{
			"abcd@efg.com",
			"1234@567.com",
		},
		ContactEmail: "test@owner.com",
	}

	tx, err := db.Begin(ctx)
	require.NoError(t, err)

	err = store.InsertSubscriptions(ctx, []*pb.Subscription{s1, s2}, tx)
	require.NoError(t, err)

	err = tx.Commit(ctx)
	require.NoError(t, err)

	actual := getSubscriptionsFromDb(t, ctx, db)
	assert.ElementsMatch(t, actual, []*pb.Subscription{s1, s2})
}

func TestInsert_EmptyList(t *testing.T) {
	ctx := context.Background()
	store, db := setUp(t)

	s := []*pb.Subscription{}

	tx, err := db.Begin(ctx)
	require.NoError(t, err)

	err = store.InsertSubscriptions(ctx, s, tx)
	require.NoError(t, err)

	err = tx.Commit(ctx)
	require.NoError(t, err)

	actual := getSubscriptionsFromDb(t, ctx, db)
	assert.Empty(t, actual)
}

func TestGet_ValidSubscription(t *testing.T) {
	ctx := context.Background()
	store, db := setUp(t)

	s := &pb.Subscription{
		Name:         "Test Subscription 1",
		Revision:     "abcd",
		BugLabels:    []string{"A", "B"},
		Hotlists:     []string{"C", "D"},
		BugComponent: "Component1>Subcomponent1",
		BugPriority:  1,
		BugSeverity:  2,
		BugCcEmails: []string{
			"abcd@efg.com",
			"1234@567.com",
		},
		ContactEmail: "test@owner.com",
	}

	insertSubscriptionToDb(t, ctx, db, s, true)
	actual, err := store.GetSubscription(ctx, "Test Subscription 1", "abcd")
	require.NoError(t, err)

	assert.Equal(t, actual, s)
}

func TestGet_AllSubscriptionsUniqByName(t *testing.T) {
	ctx := context.Background()
	store, db := setUp(t)

	s := &pb.Subscription{
		Name:         "Test Subscription 1",
		Revision:     "abcd",
		BugLabels:    []string{"A", "B"},
		Hotlists:     []string{"C", "D"},
		BugComponent: "Component1>Subcomponent1",
		BugPriority:  1,
		BugSeverity:  2,
		BugCcEmails: []string{
			"abcd@efg.com",
			"1234@567.com",
		},
		ContactEmail: "test@owner.com",
	}

	s1 := &pb.Subscription{
		Name:         "Test Subscription 2",
		Revision:     "bcde",
		BugLabels:    []string{"A", "B"},
		Hotlists:     []string{"C", "D"},
		BugComponent: "Component1>Subcomponent1",
		BugPriority:  1,
		BugSeverity:  2,
		BugCcEmails: []string{
			"abcd@efg.com",
			"1234@567.com",
		},
		ContactEmail: "test@owner.com",
	}

	insertSubscriptionToDb(t, ctx, db, s, true)
	insertSubscriptionToDb(t, ctx, db, s1, false)

	actual, err := store.GetAllSubscriptions(ctx)
	require.NoError(t, err)
	expected := []*pb.Subscription{s, s1}
	sort.Slice(actual, func(i int, j int) bool {
		return actual[i].Name < actual[j].Name
	})
	assert.Equal(t, actual, expected)
}

// Test that checks nil is returned when retrieving a non-existent subscription.
func TestGet_NonExistent(t *testing.T) {
	ctx := context.Background()
	store, _ := setUp(t)

	sub, err := store.GetSubscription(ctx, "Fake Subscription", "abcd")
	require.NoError(t, err)
	assert.Nil(t, sub)
}

// Tests that we only retrieve latest subscriptions according to
// version.
func TestGet_AllActiveSubscriptions(t *testing.T) {
	ctx := context.Background()
	store, db := setUp(t)

	s := &pb.Subscription{
		Name:         "Test Subscription 1",
		Revision:     "abcd",
		BugLabels:    []string{"A", "B"},
		Hotlists:     []string{"C", "D"},
		BugComponent: "Component1>Subcomponent1",
		BugPriority:  1,
		BugSeverity:  2,
		BugCcEmails: []string{
			"abcd@efg.com",
			"1234@567.com",
		},
		ContactEmail: "test@owner.com",
	}

	s1 := &pb.Subscription{
		Name:         "Test Subscription 2",
		Revision:     "bcde",
		BugLabels:    []string{"A", "B"},
		Hotlists:     []string{"C", "D"},
		BugComponent: "Component1>Subcomponent1",
		BugPriority:  1,
		BugSeverity:  2,
		BugCcEmails: []string{
			"abcd@efg.com",
			"1234@567.com",
		},
		ContactEmail: "test@owner.com",
	}

	insertSubscriptionToDb(t, ctx, db, s, true)
	insertSubscriptionToDb(t, ctx, db, s1, false)

	actual, err := store.GetAllActiveSubscriptions(ctx)
	require.NoError(t, err)

	expected := []*pb.Subscription{s}
	assert.Equal(t, actual, expected)
}

func insertSubscriptionToDb(t *testing.T, ctx context.Context, db pool.Pool, subscription *pb.Subscription, is_active bool) {
	const query = `INSERT INTO Subscriptions
        (name, revision, bug_labels, hotlists, bug_component, bug_priority, bug_severity, bug_cc_emails, contact_email, is_active)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	if _, err := db.Exec(ctx, query, subscription.Name, subscription.Revision, subscription.BugLabels, subscription.Hotlists, subscription.BugComponent, subscription.BugPriority, subscription.BugSeverity, subscription.BugCcEmails, subscription.ContactEmail, is_active); err != nil {
		require.NoError(t, err)
	}
}

func getSubscriptionsFromDb(t *testing.T, ctx context.Context, db pool.Pool) []*pb.Subscription {
	actual := []*pb.Subscription{}
	rows, _ := db.Query(ctx, "SELECT name, revision, bug_labels, hotlists,  bug_component, bug_priority, bug_severity, bug_cc_emails, contact_email FROM Subscriptions")
	for rows.Next() {
		subscriptionInDb := &pb.Subscription{}
		if err := rows.Scan(&subscriptionInDb.Name, &subscriptionInDb.Revision, &subscriptionInDb.BugLabels, &subscriptionInDb.Hotlists, &subscriptionInDb.BugComponent, &subscriptionInDb.BugPriority, &subscriptionInDb.BugSeverity, &subscriptionInDb.BugCcEmails, &subscriptionInDb.ContactEmail); err != nil {
			require.NoError(t, err)
		}
		actual = append(actual, subscriptionInDb)
	}
	return actual
}
