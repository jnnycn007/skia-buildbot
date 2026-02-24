import { assert } from 'chai';
import fetchMock from 'fetch-mock';
import { setUpElementUnderTest, waitForRender } from '../../../infra-sk/modules/test_util';
import { RagAppSk, Topic } from './rag-app-sk';

describe('rag-app-sk', () => {
  const newInstance = setUpElementUnderTest<RagAppSk>('rag-app-sk');

  let element: RagAppSk;

  beforeEach(async () => {
    fetchMock.get('/config', {
      instance_name: 'Test Instance',
      header_icon_url: 'test-logo.svg',
    });
    fetchMock.get('/historyrag/v1/repositories', {
      repositories: ['repo1', 'repo2'],
    });
    element = newInstance();
    await waitForRender(element);
    console.log('Element:', element);
    console.log('ShadowRoot:', element.shadowRoot);
    console.log('Is instance of RagAppSk:', element instanceof RagAppSk);
    console.log('Custom Element defined:', customElements.get('rag-app-sk'));
  });

  afterEach(() => {
    fetchMock.reset();
  });

  it('renders instance name and header icon from config', async () => {
    // Check if the instance name is rendered
    if (!element.shadowRoot) {
      throw new Error('ShadowRoot is null');
    }
    const title = element.shadowRoot.querySelector('h1');
    assert.isNotNull(title);
    assert.equal(title!.textContent, 'Test Instance');

    // Check if the header icon is rendered
    const img = element.shadowRoot!.querySelector('img');
    assert.isNotNull(img);
    assert.equal(img!.getAttribute('src'), 'test-logo.svg');
  });

  it('performs search and displays results', async () => {
    await waitForRender(element);

    const mockTopics: Topic[] = [
      { topicId: 1, topicName: 'Topic 1', summary: 'Summary 1' },
      { topicId: 2, topicName: 'Topic 2', summary: 'Summary 2' },
    ];

    fetchMock.get('/historyrag/v1/topics?query=test&topic_count=10', { topics: mockTopics });

    // Set query and trigger search
    const input = element.shadowRoot!.querySelector('md-outlined-text-field.query-input') as any;
    input.value = 'test';
    input.dispatchEvent(new Event('input'));

    const searchButton = element.shadowRoot!.querySelector('md-filled-button') as HTMLElement;
    searchButton.click();

    await waitForRender(element);

    // Verify search request was made
    assert.isTrue(fetchMock.called('/historyrag/v1/topics?query=test&topic_count=10'));

    // Verify results are displayed
    const topicItems = element.shadowRoot!.querySelectorAll('.topic-item');
    assert.equal(topicItems.length, 2);

    // Verify topic content
    const firstTopicName = topicItems[0].querySelector('.topic-name');
    assert.equal(firstTopicName!.textContent, 'Topic 1');
  });

  it('shows loading spinner when searching', async () => {
    await waitForRender(element);

    // Delay response to check loading state
    fetchMock.get(
      '/historyrag/v1/topics?query=test&topic_count=10',
      new Promise((resolve) => setTimeout(() => resolve({ topics: [] }), 100))
    );

    const input = element.shadowRoot!.querySelector('md-outlined-text-field.query-input') as any;
    input.value = 'test';
    input.dispatchEvent(new Event('input'));

    const searchButton = element.shadowRoot!.querySelector('md-filled-button') as HTMLElement;
    searchButton.click();

    await waitForRender(element);

    // Check for spinner
    const spinnerOverlay = element.shadowRoot!.querySelector('.spinner-overlay');
    assert.isNotNull(spinnerOverlay);
  });

  it('automatically generates AI summary when checkbox is checked', async () => {
    await waitForRender(element);

    const mockTopics: Topic[] = [{ topicId: 1, topicName: 'Topic 1', summary: 'Summary 1' }];
    const mockSummary = 'This is a mock AI summary.';

    fetchMock.get('/historyrag/v1/topics?query=test&topic_count=10', { topics: mockTopics });
    fetchMock.post('/historyrag/v1/summary', { summary: mockSummary });

    // Set query and trigger search
    const input = element.shadowRoot!.querySelector('md-outlined-text-field.query-input') as any;
    input.value = 'test';
    input.dispatchEvent(new Event('input'));

    // Check the AI Summary checkbox
    const checkbox = element.shadowRoot!.querySelector('md-checkbox') as any;
    checkbox.checked = true;
    checkbox.dispatchEvent(new Event('change'));

    const searchButton = element.shadowRoot!.querySelector('md-filled-button') as HTMLElement;
    searchButton.click();

    await waitForRender(element);

    // Verify both search and summary requests were made
    assert.isTrue(fetchMock.called('/historyrag/v1/topics?query=test&topic_count=10'));
    assert.isTrue(fetchMock.called('/historyrag/v1/summary'));

    // Verify summary is displayed
    const summarySection = element.shadowRoot!.querySelector('.ai-summary-content');
    assert.isNotNull(summarySection);
    assert.include(summarySection!.textContent!, 'This is a mock AI summary.');
  });

  it('performs search with a specific repository', async () => {
    await waitForRender(element);

    const mockTopics: Topic[] = [
      { topicId: 1, topicName: 'Topic 1', summary: 'Summary 1', repository: 'repo1' },
    ];

    fetchMock.get('/historyrag/v1/topics?query=test&topic_count=10&repository=repo1', {
      topics: mockTopics,
    });

    // Set query
    const input = element.shadowRoot!.querySelector('md-outlined-text-field.query-input') as any;
    input.value = 'test';
    input.dispatchEvent(new Event('input'));

    // Select repository
    const repoSelect = element.shadowRoot!.querySelector('md-outlined-select.repo-select') as any;
    repoSelect.value = 'repo1';
    repoSelect.dispatchEvent(new Event('change'));

    const searchButton = element.shadowRoot!.querySelector('md-filled-button') as HTMLElement;
    searchButton.click();

    await waitForRender(element);

    // Verify search request was made with repository parameter
    assert.isTrue(
      fetchMock.called('/historyrag/v1/topics?query=test&topic_count=10&repository=repo1')
    );

    // Verify repository name is displayed in the result
    const repoLabel = element.shadowRoot!.querySelector('.topic-repository');
    assert.isNotNull(repoLabel);
    assert.equal(repoLabel!.textContent, 'Repo: repo1');
  });

  it('passes repository when selecting a topic', async () => {
    await waitForRender(element);

    const mockTopics: Topic[] = [
      { topicId: 1, topicName: 'Topic 1', summary: 'Summary 1', repository: 'repo1' },
    ];
    const mockTopicDetails = {
      topics: [{ topicId: 1, topicName: 'Topic 1', summary: 'Full Summary', codeChunks: [] }],
    };

    fetchMock.get('/historyrag/v1/topics?query=test&topic_count=10', { topics: mockTopics });
    fetchMock.get(
      '/historyrag/v1/topic_details?topic_ids=1&include_code=true&repository=repo1&search_repository=',
      {
        topics: mockTopicDetails.topics,
      }
    );

    // Perform search
    const input = element.shadowRoot!.querySelector('md-outlined-text-field.query-input') as any;
    input.value = 'test';
    input.dispatchEvent(new Event('input'));
    const searchButton = element.shadowRoot!.querySelector('md-filled-button') as HTMLElement;
    searchButton.click();
    await waitForRender(element);

    // Click on the topic
    const topicItem = element.shadowRoot!.querySelector('.topic-item') as HTMLElement;
    topicItem.click();
    await waitForRender(element);

    // Verify topic_details request was made with repository parameter
    assert.isTrue(
      fetchMock.called(
        '/historyrag/v1/topic_details?topic_ids=1&include_code=true&repository=repo1&search_repository='
      )
    );
  });
});
