import SingleFaceGroup, { SINGLE_FACE_GROUP } from './SingleFaceGroup'
import { MemoryRouter } from 'react-router'
import 'vitest-webgl-canvas-mock'

vi.mock('../../../hooks/useScrollPagination')

import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import { MockedProvider } from "@apollo/client/testing/react";
import { ApolloError, ApolloLink } from '@apollo/client'
import { onError } from '@apollo/client/link/error'
import {
  MockedProvider,
  MockedProviderProps,
  MockedResponse
} from '@apollo/client/testing'
import { MockLink } from '@apollo/react-testing'

interface Props extends MockedProviderProps {
  mocks?: ReadonlyArray<MockedResponse>
  children?: React.ReactElement
}

const VerboseMockedProvider = (props: Props) => {
  const { mocks = [], ...otherProps } = props

  const mockLink = new MockLink(mocks)
  const errorLoggingLink = onError(
    ({ graphQLErrors, networkError }) => {
      if (graphQLErrors) {
        graphQLErrors.forEach(
          ({ message, locations, path }) =>
            console.log(
              '[GraphQL error]:' +
                `Message: ${message},` +
                `Location: ${locations},` +
                `Path: ${path}`
            )
        )
      }

      if (networkError) {
        console.log(`[Network error]: ${networkError}`)
      }
    }
  )
  const link = ApolloLink.from([errorLoggingLink, mockLink])

  return (
    <MockedProvider
      {...otherProps}
      addTypename={false}
      mocks={mocks}
      link={link}
    />
  )
}

test('single face group', async () => {
  const graphqlMocks = [
    {
      request: {
        query: SINGLE_FACE_GROUP,
        variables: { limit: 200, offset: 0, id: '123' },
      },
      result: {
        data: {
          faceGroup: {
            __typename: 'FaceGroup',
            id: '2',
            label: 'Face Group Name',
            imageFaces: [
              {
                __typename: 'ImageFace',
                id: '1',
                rectangle: {
                  __typename: 'FaceRectangle',
                  minX: 0.4912109971046448,
                  maxX: 0.5927730202674866,
                  minY: 0.2998049855232239,
                  maxY: 0.4013670086860657,
                },
                media: {
                  __typename: 'Media',
                  id: '10',
                  type: 'Photo',
                  title: '122A2785-2.jpg',
                  blurhash: 'LeIEnWD%n$Rj~pRiM{ofIpW=Iot6',
                  thumbnail: {
                    __typename: 'MediaURL',
                    url: '/photo/thumbnail_122A2785-2_jpg_lFmZcaN5.jpg',
                    width: 1024,
                    height: 1024,
                  },
                  highRes: {
                    __typename: 'MediaURL',
                    url: '/photo/122A2785-2_e4nCeMHU.jpg',
                  },
                  favorite: false,
                },
              },
              {
                __typename: 'ImageFace',
                id: '2',
                rectangle: {
                  __typename: 'FaceRectangle',
                  minX: 0.265625,
                  maxX: 0.3876950144767761,
                  minY: 0.1917019933462143,
                  maxY: 0.3705289959907532,
                },
                media: {
                  __typename: 'Media',
                  id: '52',
                  type: 'Photo',
                  title: 'image.png',
                  blurhash: 'LeIEnWD%n$Rj~pRiM{ofIpW=Iot6',
                  thumbnail: {
                    __typename: 'MediaURL',
                    url: '/photo/thumbnail_image_png_OwTDG5fM.jpg',
                    width: 1024,
                    height: 699,
                  },
                  highRes: {
                    __typename: 'MediaURL',
                    url: '/photo/image_A2YB0x3z.png',
                  },
                  favorite: false,
                },
              },
            ],
          },
        },
      },
    },
    {
      request: {
        query: SINGLE_FACE_GROUP,
        variables: { limit: 200, offset: 0 },
      },
      result: {
        data: {
          faceGroup: {
            __typename: 'FaceGroup',
            id: '2',
            label: 'Face Group Name',
            imageFaces: [
            ],
          },
        },
      },
    },
    {
      request: {
        query: SINGLE_FACE_GROUP,
        variables: { },
      },
      result: {
        data: {
          faceGroup: {
            __typename: 'FaceGroup',
            id: '2',
            label: 'Face Group Name',
            imageFaces: [
            ],
          },
        },
      },
    },
  ]

  render(
    <MemoryRouter initialEntries={['/person/123']}>
      <VerboseMockedProvider mocks={graphqlMocks}>
        <SingleFaceGroup faceGroupID="123" />
      </VerboseMockedProvider>
    </MemoryRouter>
  )

  await waitFor(() => {
    expect(screen.getAllByRole('img')).toHaveLength(2)
  })
})
