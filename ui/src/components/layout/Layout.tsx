import { gql } from '@apollo/client'
import React, { useContext } from 'react'
import { Helmet } from 'react-helmet'
import Header from '../header/Header'
import { Authorized } from '../routes/AuthorizedRoute'
import { Sidebar, SidebarContext } from '../sidebar/Sidebar'
import MainMenu from './MainMenu'

export const ADMIN_QUERY = gql`
  query adminQuery {
    myUser {
      admin
    }
  }
`

type LayoutProps = {
  children: React.ReactNode
  title: string
}

const Layout = ({ children, title, ...otherProps }: LayoutProps) => {
  const { pinned, content: sidebarContent } = useContext(SidebarContext)
  const add1 = "window.$RefreshReg$ = () => { }"
  const add2 = "window.$RefreshSig$ = () => (type) => type"

  return (
    <>
      <Helmet>
        <script type="module" src="{{ vite_dev_url }}/@vite/client"></script>
        <script type="module">
          import RefreshRuntime from 'http://localhost:5173/@react-refresh'
          RefreshRuntime.injectIntoGlobalHook(window)
          {{ add1 }}
          {{ add2 }}
          window.__vite_plugin_react_preamble_installed__ = true
        </script>
        <script type="module" src="{{ vite_dev_url }}/src/main.tsx"></script>
        <title>{title ? `${title} - Photoview` : `Photoview`}</title>
      </Helmet>
      <div className="relative" {...otherProps} data-testid="Layout">
        <Header />
        <div className="">
          <Authorized>
            <MainMenu />
          </Authorized>
          <div
            className={`mx-3 my-3 lg:mt-5 lg:mr-8 lg:ml-[292px] ${
              pinned && sidebarContent ? 'lg:pr-[420px]' : ''
            }`}
            id="layout-content"
          >
            {children}
          </div>
        </div>
        <Sidebar />
      </div>
    </>
  )
}

export default Layout
