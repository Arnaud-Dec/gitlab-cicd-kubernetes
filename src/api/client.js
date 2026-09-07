import axios from 'axios'
import { loadConfig } from '../config'
import { attachRequestNotifications } from './notifyInterceptor'

const client = attachRequestNotifications(axios.create())

client.interceptors.request.use(async (requestConfig) => {
  const { apiUrl } = await loadConfig()
  requestConfig.baseURL = apiUrl
  return requestConfig
})

export default client
