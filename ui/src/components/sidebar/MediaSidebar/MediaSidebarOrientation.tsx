import React, { useState } from 'react'
import { useMutation, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'

import { Button } from '../../../primitives/form/Input'
import { SidebarSection } from '../SidebarComponents'
import styled from 'styled-components'
import { InputLabelDescription } from '../../../Pages/SettingsPage/SettingsPage'
import { ReactComponent as LoadingSpinnerIcon } from '../../../primitives/form/icons/textboxLoadingSpinner.svg'
import { rotateMediaAction, rotateMediaActionVariables } from './__generated__/rotateMediaAction'
import { mirrorMediaAction, mirrorMediaActionVariables } from './__generated__/mirrorMediaAction'
import { removeMediaAction, removeMediaActionVariables } from './__generated__/removeMediaAction'

const ROTATE_MEDIA_MUTATION = gql`
  mutation rotateMediaAction( $mediaId: ID!) {
    rotateMedia( mediaId: $mediaId) {
      success
      message
    }
  }
`

const MIRROR_MEDIA_MUTATION = gql`
  mutation mirrorMediaAction( $mediaId: ID!) {
    mirrorMedia( mediaId: $mediaId) {
      success
      message
    }
  }
`

const REMOVE_MEDIA_MUTATION = gql`
  mutation removeMediaAction( $mediaId: ID!) {
    removeMedia( mediaId: $mediaId) {
      success
      message
    }
  }
`

type MediaSidebarOrientationProps = {
  media: MediaSidebarMedia
}

const MediaSidebarOrientation = ({ media }: MediaSidebarOrientationProps) => {
  const { t } = useTranslation()
  const [startRotateMedia, { calledRotateMedia }] = useMutation<rotateMediaAction,rotateMediaActionVariables>(ROTATE_MEDIA_MUTATION)
  const [startMirrorMedia, { calledMirrorMedia }] = useMutation<mirrorMediaAction,mirrorMediaActionVariables>(MIRROR_MEDIA_MUTATION)
  const [startRemoveMedia, { calledRemoveMedia }] = useMutation<removeMediaAction,removeMediaActionVariables>(REMOVE_MEDIA_MUTATION)

  let inputLabel = (
      <InputLabelDescription>
        {t(
          'sidebar.media.orientation.description',
          'Media Orientation'
        )}
      </InputLabelDescription>
  )

  let inputRotate = (
      <Button
        onClick={() => { startRotateMedia( { variables: { mediaId: media.id } }); }}
        disabled={calledRotateMedia}
      >
        {t('sidebar.image.rotate', 'Rotate Media')}
      </Button>
  )

  let inputMirror = (
      <Button
        onClick={() => { startMirrorMedia( { variables: { mediaId: media.id } }); }}
        disabled={calledMirrorMedia}
      >
        {t('sidebar.image.mirror', 'Mirror Media')}
      </Button>
  )

  let inputRemove = (
      <Button
        onClick={() => { startRemoveMedia( { variables: { mediaId: media.id } }); }}
        disabled={calledRemoveMedia}
      >
        {t('sidebar.image.remove', 'Remove Media')}
      </Button>
  )

  let inputSpinner = (
      <LoadingSpinnerIcon
        aria-label="Loading"
        className="inline-block animate-spin"
      />
  )

  if (calledRemoveMedia) {
    if (calledMirrorMedia) {
      if (calledRotateMedia) {
        return (
          <SidebarSection>
            { inputLabel }
            { inputRotate }
            { inputSpinner }
            { inputMirror }
            { inputSpinner }
            { inputRemove }
            { inputSpinner }
          </SidebarSection>
        )
      } else {
        return (
          <SidebarSection>
            { inputLabel }
            { inputRotate }
            { inputMirror }
            { inputSpinner }
            { inputRemove }
            { inputSpinner }
          </SidebarSection>
        )
      }
    } else {
      if (calledRotateMedia) {
        return (
          <SidebarSection>
            { inputLabel }
            { inputRotate }
            { inputSpinner }
            { inputMirror }
            { inputRemove }
            { inputSpinner }
          </SidebarSection>
        )
      } else {
        return (
          <SidebarSection>
            { inputLabel }
            { inputRotate }
            { inputMirror }
            { inputRemove }
            { inputSpinner }
          </SidebarSection>
        )
      }
    }
  } else {
    if (calledMirrorMedia) {
      if (calledRotateMedia) {
        return (
          <SidebarSection>
            { inputLabel }
            { inputRotate }
            { inputSpinner }
            { inputMirror }
            { inputSpinner }
            { inputRemove }
          </SidebarSection>
        )
      } else {
        return (
          <SidebarSection>
            { inputLabel }
            { inputRotate }
            { inputMirror }
            { inputSpinner }
            { inputRemove }
          </SidebarSection>
        )
      }
    } else {
      if (calledRotateMedia) {
        return (
          <SidebarSection>
            { inputLabel }
            { inputRotate }
            { inputSpinner }
            { inputMirror }
            { inputRemove }
          </SidebarSection>
        )
      } else {
        return (
          <SidebarSection>
            { inputLabel }
            { inputRotate }
            { inputMirror }
            { inputRemove }
          </SidebarSection>
        )
      }
    }
  }
}

export default MediaSidebarOrientation
