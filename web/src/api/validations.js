import Vue from 'vue';

export default {
  async start(validation) {
    const response = await Vue.axios.post('/api/validation-runs', validation);
    if (response.status !== 202) throw new Error('Expected 202 Accepted Status.');
    return response;
  },
  async track(id) {
    const response = await Vue.axios.get(`/api/validation-runs/${id}`);
    if (response.status !== 200) throw new Error('Expected 200 Ok Status.');
    return response;
  },
};
