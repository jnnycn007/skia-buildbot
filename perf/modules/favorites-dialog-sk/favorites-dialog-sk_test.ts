import './index';
import { expect } from 'chai';
import { FavoritesDialogSk } from './favorites-dialog-sk';

import { setUpElementUnderTest } from '../../../infra-sk/modules/test_util';
import fetchMock from 'fetch-mock';

describe('favorites-dialog-sk', () => {
  const newInstance = setUpElementUnderTest<FavoritesDialogSk>('favorites-dialog-sk');

  let element: FavoritesDialogSk;
  beforeEach(async () => {
    fetchMock.restore();
    element = newInstance((_el: FavoritesDialogSk) => {
      // Place here any code that must run after the element is instantiated but
      // before it is attached to the DOM (e.g. property setter calls,
      // document-level event listeners, etc.).
    });
    await element.updateComplete;
  });

  it('renders for new', async () => {
    expect(element).to.not.be.null;
    element.open('12345', '', '', 'url1.com');
    await element.updateComplete;

    const nameInput = element.querySelector('input[placeholder="Name"]') as HTMLInputElement;
    const urlInput = element.querySelector('input[placeholder="URL"]') as HTMLInputElement;

    expect(nameInput).to.not.be.null;
    expect(nameInput.value).to.equal('');
    expect(urlInput.value).to.equal('url1.com');
  });

  it('renders for update', async () => {
    expect(element).to.not.be.null;
    element.open('', 'Fav', 'Fav Desc', 'url.com');

    // Wait for LitElement to update the DOM
    await element.updateComplete;

    const nameInput = element.querySelector('input[placeholder="Name"]') as HTMLInputElement;
    const descInput = element.querySelector('input[placeholder="Description"]') as HTMLInputElement;
    const urlInput = element.querySelector('input[placeholder="URL"]') as HTMLInputElement;

    expect(nameInput).to.not.be.null;
    expect(descInput).to.not.be.null;
    expect(urlInput).to.not.be.null;

    expect(nameInput.value).to.equal('Fav');
    expect(descInput.value).to.equal('Fav Desc');
    expect(urlInput.value).to.equal('url.com');
  });

  it('strips begin and end by default', async () => {
    expect(element).to.not.be.null;
    element.open('', 'Fav', 'Fav Desc', 'http://localhost:8000/e/?begin=123&end=456&query=abc');

    await element.updateComplete;

    const urlInput = element.querySelector('input[placeholder="URL"]') as HTMLInputElement;
    expect(urlInput.value).to.equal('http://localhost:8000/e/?query=abc');

    const checkbox = element.querySelector('input[type="checkbox"]') as HTMLInputElement;
    expect(checkbox.checked).to.be.false;
  });

  it('includes begin and end when checkbox is checked', async () => {
    expect(element).to.not.be.null;
    element.open('', 'Fav', 'Fav Desc', 'http://localhost:8000/e/?begin=123&end=456&query=abc');

    await element.updateComplete;

    const checkbox = element.querySelector('input[type="checkbox"]') as HTMLInputElement;
    checkbox.click();
    await element.updateComplete;

    const urlInput = element.querySelector('input[placeholder="URL"]') as HTMLInputElement;
    expect(urlInput.value).to.equal('http://localhost:8000/e/?begin=123&end=456&query=abc');

    checkbox.click();
    await element.updateComplete;
    expect(urlInput.value).to.equal('http://localhost:8000/e/?query=abc');
  });

  it('saves a new favorite', async () => {
    fetchMock.post('/_/favorites/new', 200);

    // Provide the original URL with time params. By default they'll be stripped.
    const promise = element.open(
      '',
      'New Fav',
      'A description',
      'http://localhost/test?begin=1&end=2'
    );
    await element.updateComplete;

    const saveBtn = element.querySelectorAll('.buttons button')[1] as HTMLButtonElement;
    saveBtn.click();

    await promise;

    expect(fetchMock.called('/_/favorites/new')).to.be.true;
    const requestArgs = fetchMock.lastCall('/_/favorites/new');
    expect(requestArgs).to.not.be.undefined;
    const requestBody = JSON.parse(requestArgs![1]!.body!.toString());

    expect(requestBody).to.deep.equal({
      name: 'New Fav',
      description: 'A description',
      url: 'http://localhost/test',
    });
  });

  it('saves an existing favorite', async () => {
    fetchMock.post('/_/favorites/edit', 200);

    const promise = element.open('123', 'Existing Fav', 'Desc', 'http://localhost/test');
    await element.updateComplete;

    const saveBtn = element.querySelectorAll('.buttons button')[1] as HTMLButtonElement;
    saveBtn.click();

    await promise;

    expect(fetchMock.called('/_/favorites/edit')).to.be.true;
    const requestArgs = fetchMock.lastCall('/_/favorites/edit');
    expect(requestArgs).to.not.be.undefined;
    const requestBody = JSON.parse(requestArgs![1]!.body!.toString());

    expect(requestBody).to.deep.equal({
      id: '123',
      name: 'Existing Fav',
      description: 'Desc',
      url: 'http://localhost/test',
    });
  });
});
