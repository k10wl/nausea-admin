const STORAGE_KEY = "nausea_admin_prefferences";

const DEFAULT_VALUES = {
  version: "0.1.0",
  showDeleted: true,
};

class NauseaStore {
  constructor(instance) {
    /** @type {Storage} */
    this.$instance = instance;
  }

  get(key) {
    const fromDefaultsWithKey = this.$fromDefaults.bind(null, key);
    try {
      const state = this.$getState();
      if (state === null) {
        return fromDefaultsWithKey();
      }
      const decoded = this.$decode(state);
      return decoded[key] ?? decoded ?? fromDefaultsWithKey();
    } catch (e) {
      return fromDefaultsWithKey();
    }
  }

  set(key, value) {
    let state = this.$getState();
    if (state === null) {
      state = DEFAULT_VALUES;
    } else {
      state = this.$decode(state);
    }
    state[key] = value;
    this.$setState(state);
  }

  $fromDefaults(key) {
    return DEFAULT_VALUES[key] ?? DEFAULT_VALUES;
  }

  $getState() {
    return this.$instance.getItem(STORAGE_KEY);
  }

  $setState(state) {
    return this.$instance.setItem(STORAGE_KEY, this.$encode(state));
  }

  $decode(encoded) {
    try {
      return JSON.parse(encoded);
    } catch {
      return DEFAULT_VALUES;
    }
  }

  $encode(decoded) {
    return JSON.stringify(decoded);
  }
}

const $ = new NauseaStore(localStorage);
