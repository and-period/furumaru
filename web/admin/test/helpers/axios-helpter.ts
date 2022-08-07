import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'

const axiosMock = new MockAdapter(axios)
const baseURL = process.env.API_BASE_URL || 'http://localhost:18010'

export { axiosMock, baseURL }
