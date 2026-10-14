import React, { useEffect } from 'react'
import { Helmet } from 'react-helmet'
import Routes from './components/routes/Routes'
import Messages from './components/messages/Messages'
import { useTranslation } from 'react-i18next'
import { loadTranslations } from './localization'
import { useLocation } from 'react-router'

const App = () => {
  const { t } = useTranslation()
  const { pathname } = useLocation()
  const add1 = "window.$RefreshReg$ = () => { }"
  const add2 = "window.$RefreshSig$ = () => (type) => type"

  loadTranslations()

  useEffect(() => {
    window.scrollTo(0, 0)
    if (document.activeElement != document.body)
      (document.activeElement as HTMLInputElement).blur()
  }, [pathname])

  return (
    <>
      <Helmet>
        <meta
          name="description"
          content={t(
            'meta.description',
            'Simple and User-friendly Photo Gallery for Personal Servers'
          )}
        />
        <script type="module" src="{{ vite_dev_url }}/@vite/client"></script>
        <script type="module">
          import RefreshRuntime from 'http://localhost:1234/@react-refresh'
          RefreshRuntime.injectIntoGlobalHook(window)
          {{ add1 }}
          {{ add2 }}
          window.__vite_plugin_react_preamble_installed__ = true
        </script>
        <script type="module" src="{{ vite_dev_url }}/src/main.tsx"></script>
      </Helmet>
      <Routes />
      <Messages />
    </>
  )
}

export default App
