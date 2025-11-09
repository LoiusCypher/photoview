import { render, screen, waitFor } from '@testing-library/react'
import React from 'react'
import LoginPage from './LoginPage'
import * as authentication from '../../helpers/authentication'
import { createMemoryRouter, BrowserRouter, RouterProvider } from "react-router"
import { MockedProvider } from "@apollo/client/testing/react";
import { mockInitialSetupGraphql } from './loginTestHelpers'

vi.mock('../../helpers/authentication.ts')

const authToken = vi.mocked(authentication.authToken)

  const routes = [
    {
      path: '/',
      element: <LoginPage />,
      //Component: <>App </>,
      children: [
        {
          path: 'login',
          Component: <LoginPage />,
        },
        {
          path: 'initialSetup',
          element: <>InitialSetupPage </>,
        },
      ],
    },
  ];

  const routerRoot = createMemoryRouter(
    routes,
    {
      // Set for where you want to start in the routes. Remember, KISS (Keep it simple, stupid) the routes.
      initialEntries: ['/'],
      // We don't need to explicitly set this, but it's nice to have.
      initialIndex: 0,
      errorElement: <>NotFound </>,
    },
  )

  const routerLogin = createMemoryRouter(
    routes,
    {
      // Set for where you want to start in the routes. Remember, KISS (Keep it simple, stupid) the routes.
      initialEntries: ['/login'],
      // We don't need to explicitly set this, but it's nice to have.
      initialIndex: 0,
      errorElement: <>NotFound </>,
    },
  )

describe('Login page redirects', () => {
  test('Auth token redirect', async () => {
    authToken.mockImplementation(() => 'some-token')

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(false)]}>
        <RouterProvider router={routerLogin} />
      </MockedProvider>
    )

    await waitFor(() => {
      expect(routerLogin.state.location.pathname).toBe('/')
    })
  })

  test('Initial setup redirect', async () => {
    authToken.mockImplementation(() => null)

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(true)]}>
        <RouterProvider router={routerLogin}>
          <LoginPage />
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(routerLogin.state.location.pathname).toBe('/initialSetup')
    })
  })
})

describe('Login page', () => {
  test('Render login form', () => {
    authToken.mockImplementation(() => null)

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(false)]}>
        <RouterProvider router={routerLogin}>
          <LoginPage />
        </RouterProvider>
      </MockedProvider>
    )

    expect(screen.getByLabelText('Username')).toBeInTheDocument()
    expect(screen.getByLabelText('Password')).toBeInTheDocument()
    expect(screen.getByDisplayValue('Sign in')).toBeInTheDocument()
  })
})
