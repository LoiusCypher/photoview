import { render, screen, waitFor } from '@testing-library/react'
import React from 'react'
import LoginPage from './LoginPage'
import * as authentication from '../../helpers/authentication'
import { createMemoryRouter, BrowserRouter, RouterProvider } from "react-router-dom"
import { MockedProvider } from "@apollo/client/testing/react";
import { mockInitialSetupGraphql } from './loginTestHelpers'

vi.mock('../../helpers/authentication.ts')

const authToken = vi.mocked(authentication.authToken)

//const setupMyTest = (url: string) => {
  const router = createMemoryRouter(
    [
      {
        path: '/',
        element: <>App </>,
      },
      {
        path: '/login',
        element: <>LoginPage </>,
      },
    ],
    {
      // Set for where you want to start in the routes. Remember, KISS (Keep it simple, stupid) the routes.
      initialEntries: ['/'],
      // We don't need to explicitly set this, but it's nice to have.
      initialIndex: 0,
    }
  )

  //render(<RouterProvider router={router} />)

  // Objectify the router so we can explicitly pull when calling setupMyTest
  //return { router }
//}

  const router2 = createMemoryRouter(
    [
      {
        path: '/',
        element: <>App </>,
      },
      {
        path: '/login',
        element: <>Login </>,
      },
      {
        path: '/initialSetup',
        element: <>InitialSetupPage </>,
      },
    ],
    {
      // Set for where you want to start in the routes. Remember, KISS (Keep it simple, stupid) the routes.
      initialEntries: ['/login'],
      // We don't need to explicitly set this, but it's nice to have.
      initialIndex: 0,
    }
  )

describe('Login page redirects', () => {
  test('Auth token redirect', async () => {
    authToken.mockImplementation(() => 'some-token')

    //const { router } = setupMyTest('/')

    render(
      <MockedProvider mocks={[]}>
        <RouterProvider router={router}>
          <LoginPage />
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(router.state.location.pathname).toBe('/')
    })
  })

  test('Initial setup redirect', async () => {
    authToken.mockImplementation(() => null)

    //const { router } = setupMyTest('/login')

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(true)]}>
        <RouterProvider router={router2}>
          <LoginPage />
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(router.state.location.pathname).toBe('/initialSetup')
    })
  })
})

describe('Login page', () => {
  test('Render login form', () => {
    authToken.mockImplementation(() => null)

    //const { router } = setupMyTest('/login')

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(false)]}>
        <RouterProvider router={router2}>
          <LoginPage />
        </RouterProvider>
      </MockedProvider>
    )

    expect(screen.getByLabelText('Username')).toBeInTheDocument()
    expect(screen.getByLabelText('Password')).toBeInTheDocument()
    expect(screen.getByDisplayValue('Sign in')).toBeInTheDocument()
  })
})
