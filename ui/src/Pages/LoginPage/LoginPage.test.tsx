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
    element: <LoginPage >,
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

describe('Login page redirects', () => {
  test('Auth token redirect', async () => {
    authToken.mockImplementation(() => 'some-token')

    const history = createMemoryRouter(routes, {
      initialEntries: ['/login'],
    })

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(false)]}>
        <RouterProvider router={history}>
          <LoginPage />
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(history.state.location.pathname).toBe('/')
    })
  })

  test('Initial setup redirect', async () => {
    authToken.mockImplementation(() => null)

    const history = createMemoryRouter(routes, {
      initialEntries: ['/login'],
    })

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(true)]}>
        <RouterProvider router={history}>
          <LoginPage />
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(history.state.location.pathname).toBe('/initialSetup')
    })
  })
})

describe('Login page', () => {
  test('Render login form', () => {
    authToken.mockImplementation(() => null)

    const history = createMemoryRouter(routes, {
      initialEntries: ['/login'],
    })

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(false)]}>
        <RouterProvider router={history}>
          <LoginPage />
        </RouterProvider>
      </MockedProvider>
    )

    expect(screen.getByLabelText('Username')).toBeInTheDocument()
    expect(screen.getByLabelText('Password')).toBeInTheDocument()
    expect(screen.getByDisplayValue('Sign in')).toBeInTheDocument()
  })
})
