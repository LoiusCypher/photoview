import React, { useState } from 'react'
import { useMutation, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'

import { Button } from '../../../primitives/form/Input'
import { SidebarSection } from '../SidebarComponents'

//import { SIDEBAR_MEDIA_QUERY } from './MediaSidebar'
import { sidebarMediaQuery } from './__generated__/sidebarMediaQuery'

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

type MediaSidebarMediaFacesProps = {
  faces: sidebarMediaQuery_media_faces
}

const MediaSidebarRescan = ({ media }: MediaSidebarFacesProps) => {
  const { t } = useTranslation()
  const [startMediaScanner, { calledMediaScan }] = useMutation<scanMediaAction,scanMediaActionVariables>(SCAN_MEDIA_MUTATION)
  const [startReMediaScanner, { calledReMedia }] = useMutation<reScanMediaAction,reScanMediaActionVariables>
							(RE_SCAN_MEDIA_MUTATION, {
							  refetchQueries: [
//							    SIDEBAR_MEDIA_QUERY, // DocumentNode object parsed with gql
							    "sidebarMediaQuery", // Query name
							  ],
							});

    let reScan = (
      <Button
        onClick={() => { startReMediaScanner( { variables: { mediaId: media.id } }); }}
        disabled={calledReMedia}
      >
        {t('sidebar.people.rescan.original', 'Rescan Media')}
      </Button>
    )

    let mediaScan = (
      <Button
        onClick={() => { startMediaScanner( { variables: { mediaId: media.id } }); }}
        disabled={calledMediaScan}
      >
        {t('sidebar.people.rescan.thumbnail', 'Scan album containig thumbnail')}
      </Button>
    )

    if (calledMediaScan) 
      mediaScan = (
        <div className="relative">
          {mediaScan}
          <LoadingSpinnerIcon
            aria-label="Loading"
            className="absolute right-[8px] top-[7px] animate-spin"
          />
        </div>
      )

    if (calledReMedia)
      reScan = (
        <div className="relative">
          {reScan}
          <LoadingSpinnerIcon
            aria-label="Loading"
            className="absolute right-[8px] top-[7px] animate-spin"
          />
        </div>
    )

      return (
        <SidebarSection>
          <InputLabelDescription>
            {t(
              'sidebar.people.rescan.description',
              'Gesichtserkennung'
            )}
          </InputLabelDescription>
          { mediaScan }
          { reScan }
        </SidebarSection>
      )
}

export default MediaSidebarRescan
