const SANITIZE_KEYS = ['ix', 'mn', 'cix', 'np', 'l', 'td', 'cl', 'ct'];
const SANITIZE_KEYS_IF_ROOT = ['v', 'props'];
export const SANITIZE_KEYS_IF_SLOT = ['t'];
const SANITIZE_KEYS_IF_NOT_LAYER = ['ind'];
const SANITIZE_KEYS_IF_ASSET = ['fr'];
const SANITIZE_KEYS_IF_NOT_PRECOMP = ['sr', 'st'];
const SANITIZE_KEYS_IF_VAL_IS_ZERO = ['ddd', 'bm', 'ao'];

const PRECOMP_LAYER_TY = 0;

export const COMMON_EXPORTER_FIELDS = [
  ...SANITIZE_KEYS,
  ...SANITIZE_KEYS_IF_ASSET,
  ...SANITIZE_KEYS_IF_NOT_LAYER,
  ...SANITIZE_KEYS_IF_NOT_PRECOMP,
  ...SANITIZE_KEYS_IF_ROOT,
  ...SANITIZE_KEYS_IF_VAL_IS_ZERO,
];

/**
 * Sanitize lottie removes fields that are not part of the spec but
 * are added by common exporters. This function will mutate the lottie,
 * so the caller should provide a copy if that behavior is undesireable.
 * @param lottie
 */
export function sanitizeLottie(lottie: any) {
  for (const key of SANITIZE_KEYS_IF_ROOT) {
    delete lottie[key];
  }

  sanitizeNode(lottie);

  if (lottie.slots) {
    sanitizeSlots(lottie);
  }
}

function sanitizeSlots(lottie: any) {
  Object.values(lottie.slots).forEach((slot: any) => {
    for (const key of SANITIZE_KEYS_IF_SLOT) {
      delete slot[key];
    }
  });
}

function sanitizeNode(node: any) {
  if (Array.isArray(node)) {
    node.forEach((subNode) => sanitizeNode(subNode));
  } else if (typeof node === 'object') {
    for (const key of SANITIZE_KEYS) {
      delete node[key];
    }

    for (const key of SANITIZE_KEYS_IF_VAL_IS_ZERO) {
      if (node[key] === 0) {
        delete node[key];
      }
    }

    if (typeof node['ty'] === 'number') {
      // Layer object
      if (node['ty'] !== PRECOMP_LAYER_TY) {
        for (const key of SANITIZE_KEYS_IF_NOT_PRECOMP) {
          delete node[key];
        }
      }
    } else {
      for (const key of SANITIZE_KEYS_IF_NOT_LAYER) {
        delete node[key];
      }
    }

    if (typeof node['id'] === 'string') {
      // Asset
      for (const key of SANITIZE_KEYS_IF_ASSET) {
        delete node[key];
      }
    }

    // Recurse
    Object.values(node).forEach((subNode) => sanitizeNode(subNode));
  }
}
