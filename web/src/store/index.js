import Vue from 'vue';
import Vuex from 'vuex';
import createLogger from 'vuex/dist/logger';

import config from './modules/config';

Vue.use(Vuex);

const strict = process.env.NODE_ENV !== 'production';
const plugins = process.env.NODE_ENV !== 'production' ? [createLogger()] : [];

// to debug the store install the Vue.js chrome/firefox extension
export default new Vuex.Store({
  modules: {
    config,
  },
  strict,
  plugins,
});
