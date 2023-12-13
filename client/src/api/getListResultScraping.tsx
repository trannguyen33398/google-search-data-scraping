

import { ScrapingDatas } from '../type/scrapingItem.type'
import http from '../util/http'

export const getListResultScraping = (page: number | string, limit: number | string, signal?: AbortSignal) =>
  http.get<ScrapingDatas>('/history-upload', {
    params: {
      _page: page,
      _limit: limit
    },
    signal
  })

  export const uploadFile = ( form: FormData,signal?: AbortSignal) =>
  http.post<ScrapingDatas>('/upload', {
    form,
    signal
  })