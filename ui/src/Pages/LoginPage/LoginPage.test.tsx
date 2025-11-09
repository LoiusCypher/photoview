import { render, screen, waitFor } from '@testing-library/react'
import React from 'react'
import LoginPage from './LoginPage'
import * as authentication from '../../helpers/authentication'
import { createMemoryRouter, BrowserRouter, RouterProvider } from "react-router"
import { MockedProvider } from "@apollo/client/testing/react";
import { mockInitialSetupGraphql } from './loginTestHelpers'

vi.mock('../../helpers/authentication.ts')

const authToken = vi.mocked(authentication.authToken)

  const FAKE_EVENT_2 = { name: "test event 2" };
  const FAKE_EVENT_3 = { name: "test event 3" };
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

  const router1 = createMemoryRouter(
    routes,
    {
      // Set for where you want to start in the routes. Remember, KISS (Keep it simple, stupid) the routes.
      initialEntries: ['/'],
      // We don't need to explicitly set this, but it's nice to have.
      initialIndex: 0,
      errorElement: <>NotFound </>,
    },
  )

  const router2 = createMemoryRouter(
    routes,
    {
      // Set for where you want to start in the routes. Remember, KISS (Keep it simple, stupid) the routes.
      //initialEntries: ['/login'],
      initialEntries: ['/'],
      // We don't need to explicitly set this, but it's nice to have.
      initialIndex: 0,
      errorElement: <>NotFound </>,
    },
  )

  const router3 = createMemoryRouter(
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
      <MockedProvider mocks={[]}>
        <RouterProvider router={router1}>
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(router1.state.location.pathname).toBe('/')
    })
  })

  test('Initial setup redirect', async () => {
    authToken.mockImplementation(() => null)

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(true)]}>
        <RouterProvider router={router2}>
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(router2.state.location.pathname).toBe('/initialSetup')
    })
  })
})

describe('Login page', () => {
  test('Render login form', async () => {
    authToken.mockImplementation(() => null)

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(false)]}>
        <RouterProvider router={router3}>
          <LoginPage />
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(router3.state.location.pathname).toBe('/login')
    })
    expect(screen.getByLabelText('Username')).toBeInTheDocument()
    expect(screen.getByLabelText('Password')).toBeInTheDocument()
    expect(screen.getByDisplayValue('Sign in')).toBeInTheDocument()
  })
})
