/**
 * @module modules/favorites-dialog-sk
 * @description <h2><code>favorites-dialog-sk</code></h2>
 *
 * This module is a modal that contains a form to capture user
 * input for adding/editing a new favorite.
 */
import { html, LitElement } from 'lit';
import { customElement, property, query, state } from 'lit/decorators.js';
import { errorMessage } from '../../../elements-sk/modules/errorMessage';
import '../../../elements-sk/modules/spinner-sk';
import '../../../elements-sk/modules/icons/close-icon-sk';

// FavoritesDialogSk is a modal that contains a form to capture user
// input for adding/editing a new favorite.
@customElement('favorites-dialog-sk')
export class FavoritesDialogSk extends LitElement {
  private static nextUniqueId = 0;

  private readonly uniqueId = `${FavoritesDialogSk.nextUniqueId++}`;

  @property({ type: String })
  favId: string = '';

  @property({ type: String })
  name: string = '';

  @property({ type: String })
  description: string = '';

  @property({ type: String })
  url: string = '';

  @query('dialog')
  private dialog!: HTMLDialogElement;

  @state()
  private updatingFavorite: boolean = false;

  private resolve: ((value?: any) => void) | null = null;

  private reject: ((value?: any) => void) | null = null;

  createRenderRoot() {
    return this;
  }

  private dismiss(): void {
    this.dialog.close();
    if (this.reject) {
      this.reject();
    }
  }

  private async confirm(): Promise<void> {
    if (this.name === '' || this.url === '') {
      errorMessage('Name and url must be non empty');
      return;
    }

    let apiUrl = '/_/favorites/new';
    let body: {
      id?: string;
      name: string;
      description: string;
      url: string;
    } = {
      name: this.name,
      description: this.description,
      url: this.url,
    };
    if (this.favId !== '') {
      body = { ...body, id: this.favId };
      apiUrl = '/_/favorites/edit';
    }

    try {
      this.updatingFavorite = true;
      const resp = await fetch(apiUrl, {
        method: 'POST',
        body: JSON.stringify(body),
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!resp.ok) {
        const msg = await resp.text();
        errorMessage(`${resp.statusText}: ${msg}`);
      }
    } finally {
      this.updatingFavorite = false;
      this.dialog.close();
      if (this.resolve) {
        this.resolve();
      }
    }
  }

  // open shows the popup dialog when called.
  public async open(
    favId?: string,
    name?: string,
    description?: string,
    url?: string
  ): Promise<void> {
    this.favId = favId || '';
    this.name = name || '';
    this.description = description || '';
    this.url = url || window.location.href;

    await this.updateComplete;
    this.dialog.showModal();

    // If the dialog closes it could be due to 2 reasons:
    //    1: User pressed on close
    //    2: The favorite got added/edited.
    // In this module, we want to re-fetch the favorites when the dialog is closed
    // but we only want to re-fetch if closed due to reason 2.
    // So we're using the reject function when the user presses on close dialog
    // which is eventually used in favorites-sk to decide if it wants to
    // re-fetch the favorites or not.
    return await new Promise((resolve, reject) => {
      this.resolve = resolve;
      this.reject = reject;
    });
  }

  render() {
    return html`
      <dialog class="fav-dialog">
        <h2>Favorite</h2>

        <button class="close-btn" @click=${this.dismiss}>
          <close-icon-sk></close-icon-sk>
        </button>

        <div class="spin-container">
          <spinner-sk ?active=${this.updatingFavorite}></spinner-sk>
        </div>

        <span class="label">
          <label for="name-${this.uniqueId}">Name*</label>
        </span>
        <input
          id="name-${this.uniqueId}"
          placeholder="Name"
          .value="${this.name}"
          @input=${(e: Event) => (this.name = (e.target as HTMLInputElement).value)}>
        </input>
        <br/>

        <span class="label">
          <label for="desc-${this.uniqueId}">Description</label>
        </span>
        <input
          id="desc-${this.uniqueId}"
          placeholder="Description"
          .value="${this.description}"
          @input=${(e: Event) => (this.description = (e.target as HTMLInputElement).value)}></input>
        <br/>

        <span class="label">
          <label for="url-${this.uniqueId}">URL*</label>
        </span>
        <input
          id="url-${this.uniqueId}"
          placeholder="URL"
          .value="${this.url}"
          @input=${(e: Event) => (this.url = (e.target as HTMLInputElement).value)}></input>
        <br/><br/>

        <div ?hidden="${!this.updatingFavorite}">
          Working on it...
        </div>

        <div class="buttons">
          <button ?disabled="${this.updatingFavorite}" @click=${this.dismiss}>Cancel</button>
          <button ?disabled="${this.updatingFavorite}" @click=${this.confirm}>Save</button>
        </div>
      </dialog>`;
  }
}
