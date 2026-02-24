import fetchMock from 'fetch-mock';

fetchMock.get('glob:/historyrag/v1/repositories', {
  repositories: ['skia', 'infra', 'chrome'],
});

fetchMock.get('glob:/historyrag/v1/topics?*', (url) => {
  const params = new URLSearchParams(url.split('?')[1]);
  const repo = params.get('repository');
  return {
    topics: [
      {
        topicId: 1,
        topicName: 'Database Migration',
        summary: `Discussion about migrating from MySQL to Spanner.${
          repo ? ` (Repo: ${repo})` : ''
        }`,
        repository: repo || 'skia',
      },
      {
        topicId: 2,
        topicName: 'Authentication Bug',
        summary: `Fixing the token expiration issue in auth service.${
          repo ? ` (Repo: ${repo})` : ''
        }`,
        repository: repo || 'infra',
      },
      {
        topicId: 3,
        topicName: 'UI Refactoring',
        summary: `Updating the dashboard to use Material 3.${repo ? ` (Repo: ${repo})` : ''}`,
        repository: repo || 'infra',
      },
    ],
  };
});

fetchMock.get('glob:/historyrag/v1/topic_details?*', {
  topics: [
    {
      topicId: 1,
      topicName: 'Database Migration',
      summary:
        'Discussion about migrating from MySQL to Spanner. The migration involves schema changes and data backfill.',
      codeChunks: ['func Migrate() {\n  // Connect to Spanner\n}', 'CREATE TABLE users (...)'],
    },
  ],
});

import './index';
