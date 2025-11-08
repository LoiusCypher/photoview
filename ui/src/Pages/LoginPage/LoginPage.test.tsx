import { render, screen, waitFor } from '@testing-library/react'
import React from 'react'
import LoginPage from './LoginPage'
import * as authentication from '../../helpers/authentication'
import { createMemoryRouter, BrowserRouter, RouterProvider } from "react-router-dom"
import { MockedProvider } from "@apollo/client/testing/react";
import { mockInitialSetupGraphql } from './loginTestHelpers'

vi.mock('../../helpers/authentication.ts')

const authToken = vi.mocked(authentication.authToken)

  const FAKE_EVENT_2 = { name: "test event 2" };
  const routes = [
    {
      path: '/',
      element: <>App </>,
    },
    {
      path: '/login',
      element: <LoginPage />,
      loader: () => FAKE_EVENT_2,
    },
    {
      path: '/initialSetup',
      element: <>InitialSetupPage </>,
    },
  ];

  const router = createMemoryRouter(
    routes,
    {
      // Set for where you want to start in the routes. Remember, KISS (Keep it simple, stupid) the routes.
      initialEntries: ['/', '/login'],
      // We don't need to explicitly set this, but it's nice to have.
      initialIndex: 0,
      errorElement: <>NotFound </>,
//      children: [
//        { path: "/", element: <Home /> },
//        { path: "/posts", element: <Posts /> },
//        { path: "/post/:postId", element: <PostPage /> },
//      ],
    },
  )

//  const router2 = createMemoryRouter(
//    [
//      {
//        path: '/',
//        element: <>App </>,
//      },
//      {
//        path: '/login',
//        element: <LoginPage />,
//      },
//      {
//        path: '/initialSetup',
//        element: <>InitialSetupPage </>,
//      },
//    ],
//    {
//      // Set for where you want to start in the routes. Remember, KISS (Keep it simple, stupid) the routes.
//      initialEntries: ['/login'],
//      // We don't need to explicitly set this, but it's nice to have.
//      initialIndex: 0,
//    }
//  )

describe('Login page redirects', () => {
  test('Auth token redirect', async () => {
    authToken.mockImplementation(() => 'some-token')

    //const { router } = setupMyTest('/')

    render(
      <MockedProvider mocks={[]}>
        <RouterProvider router={router}>
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
          //<LoginPage />

    render(
      <MockedProvider mocks={[mockInitialSetupGraphql(true)]}>
        <RouterProvider router={router}>
        </RouterProvider>
      </MockedProvider>
    )

    await waitFor(() => {
      expect(router.state.location.pathname).toBe('/initialSetup')
    })
  })


//describe('Login page', () => {
//  test('Render login form', () => {
//    authToken.mockImplementation(() => null)
//
//    //const { router } = setupMyTest('/login')
//
//    render(
//      <MockedProvider mocks={[mockInitialSetupGraphql(false)]}>
//        <RouterProvider router={router2}>
//          <LoginPage />
//        </RouterProvider>
//      </MockedProvider>
//    )

//    expect(screen.getByLabelText('Username')).toBeInTheDocument()
//    expect(screen.getByLabelText('Password')).toBeInTheDocument()
//    expect(screen.getByDisplayValue('Sign in')).toBeInTheDocument()
//  })
})
