import axios, { AxiosInstance } from 'axios'

class Http {
  instance: AxiosInstance
  constructor() {
    this.instance = axios.create({
      baseURL: 'http://localhost:9090',
      timeout: 10000,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}

const http = new Http().instance

export default http