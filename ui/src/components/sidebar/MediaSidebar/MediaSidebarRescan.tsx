import React, { useState } from 'react'
import { useMutation, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'

import { Button } from '../../../primitives/form/Input'
import { SidebarSection } from '../SidebarComponents'
import styled from 'styled-components'
import { scanMediaAction, scanMediaActionVariables } from './__generated__/scanMediaAction'
import { reScanMediaAction, reScanMediaActionVariables } from './__generated__/reScanMediaAction'
import { InputLabelDescription } from '../../../Pages/SettingsPage/SettingsPage'
import { ReactComponent as LoadingSpinnerIcon } from '../../../primitives/form/icons/textboxLoadingSpinner.svg'

const SCAN_MEDIA_MUTATION = gql`
  mutation scanMediaAction( $mediaId: ID!) {
    scanMedia( mediaId: $mediaId) {
      success
      message
    }
  }
`

const RE_SCAN_MEDIA_MUTATION = gql`
  mutation reScanMediaAction( $mediaId: ID!) {
    reScanMedia( mediaId: $mediaId) {
      success
      message
    }
  }
`

type MediaSidebarFacesProps = {
  media: MediaSidebarMedia
}

const MediaSidebarRescan = ({ media }: MediaSidebarFacesProps) => {
  const { t } = useTranslation()
  const [startMediaScanner, { calledMediaScan }] = useMutation<scanMediaAction,scanMediaActionVariables>(SCAN_MEDIA_MUTATION)
  const [startReMediaScanner, { calledReMedia }] = useMutation<reScanMediaAction,reScanMediaActionVariables>(RE_SCAN_MEDIA_MUTATION)

  let inputLabel = (
      <InputLabelDescription>
        {t(
          'sidebar.people.rescan.description',
          'Gesichtserkennung'
        )}
      </InputLabelDescription>
  )

  let inputScan = (
      <Button
        onClick={() => { startMediaScanner( { variables: { mediaId: media.id } }); }}
        disabled={calledMediaScan}
      >
        {t('sidebar.people.rescan.thumbnail', 'Scan album containig thumbnail')}
      </Button>
  )

  let inputReScan = (
      <Button
        onClick={() => { startReMediaScanner( { variables: { mediaId: media.id } }); }}
        disabled={calledReMedia}
      >
        {t('sidebar.people.rescan.thumbnail', 'Rescan Thumbnail')}
      </Button>
  )

  let inputSpinner = (
      <LoadingSpinnerIcon
        aria-label="Loading"
        className="inline-block animate-spin"
      />
  )

  if (calledReMedia) {
    if (calledMediaScan) {
      return (
        <SidebarSection>
          { inputLabel }
          { inputScan }
          { inputSpinner }
          { inputReScan }
        </SidebarSection>
      )
    } else {
      return (
        <SidebarSection>
          { inputLabel }
          { inputScan }
          { inputReScan }
        </SidebarSection>
      )
    }
  } else {
    if (calledMediaScan) {
      return (
        <SidebarSection>
          { inputLabel }
          { inputScan }
          { inputSpinner }
          { inputReScan }
          { inputSpinner }
        </SidebarSection>
      )
    } else {
      return (
        <SidebarSection>
          { inputLabel }
          { inputScan }
          { inputReScan }
          { inputSpinner }
        </SidebarSection>
      )
    }
  }
}

export default MediaSidebarRescan
