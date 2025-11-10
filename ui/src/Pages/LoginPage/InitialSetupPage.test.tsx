import React from 'react'
import { MockedProvider } from "@apollo/client/testing/react";
import { render, screen, waitFor } from '@testing-library/react'
import { createMemoryRouter, BrowserRouter, RouterProvider } from "react-router"
import * as authentication from '../../helpers/authentication'
import InitialSetupPage from './InitialSetupPage'
import { mockInitialSetupGraphql } from './loginTestHelpers'

vi.mock('../../helpers/authentication.ts')

const authToken = vi.mocked(authentication.authToken)

  const router2 = createMemoryRouter(
    [
      {
        path: '/',
        Component: <InitialSetupPage />,
        children: [
          {
            path: 'initialSetup',
          Component: <InitialSetupPage />,
          },
        ],
      },
    ],
    {
      // Set for where you want to start in the routes. Remember, KISS (Keep it simple, stupid) the routes.
      initialEntries: ['/initialSetup'],
      // We don't need to explicitly set this, but it's nice to have.
      initialIndex: 0,
      errorElement: <>NotFound </>,
    },
  )

  const routes = [
    {
      path: '/',
      element: <>Root </>,
      children: [
        {
          path: 'initialSetup',
          Component: <InitialSetupPage />,
        },
      ],
    },
  ];

  const router1 = createMemoryRouter(
    routes,
    {
      // Set for where you want to start in the routes. Remember, KISS (Keep it simple, stupid) the routes.
      initialEntries: ['/initialSetup'],
      // We don't need to explicitly set this, but it's nice to have.
      initialIndex: 0,
      errorElement: <>NotFound </>,
    },
  )

describe('Initial setup page', () => {
  test('Render initial setup form', () => {
    authToken.mockImplementation(() => null)

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(true)]}>
        <RouterProvider router={router2}>
          <InitialSetupPage />
        </RouterProvider>
      </MockedProvider>
    )

    expect(screen.getByLabelText('Username')).toBeInTheDocument()
    expect(screen.getByLabelText('Password')).toBeInTheDocument()
    expect(screen.getByLabelText('Photo path')).toBeInTheDocument()
    expect(screen.getByDisplayValue('Setup Photoview')).toBeInTheDocument()
  })

  test('Redirect if auth token is present', async () => {
    authToken.mockImplementation(() => 'some-token')

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(true)]}>
        <RouterProvider router={router2}>
          <InitialSetupPage />
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(router2.state.location.pathname).toBe('/')
    })
  })

  test('Redirect if not initial setup', async () => {
    authToken.mockImplementation(() => null)

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(false)]}>
        <RouterProvider router={router1}>
          <InitialSetupPage />
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(router1.state.location.pathname).toBe('/')
    })
  })
})
