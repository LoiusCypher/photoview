import { useMutation, gql } from '@apollo/client'
import React, { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'
import { Button } from '../../primitives/form/Input'
import {
  InputLabelDescription,
  InputLabelTitle,
  SectionTitle,
} from './SettingsPage'
import { exportFaces } from './__generated__/exportFaces'

const EXPORT_ALL_FACES = gql`
  mutation exportFaces {
    exportAllFaces {
      success
      message
    }
  }
`

const ExportAllFacesButton = () => {
  const { t } = useTranslation()
  const [startExport, { called }] = useMutation<exportFaces>(EXPORT_ALL_FACES)

  return (
    <Button
      className="mb-4"
      onClick={() => {
        startExport()
      }}
      disabled={called}
    >
      {t('settings.export_all_faces', 'Export All Faces')}
    </Button>
  )
}

const DeveloperToolsWrapper = styled.div`
  margin-bottom: 24px;
`

const DeveloperTools = () => {
  const { t } = useTranslation()
  const [startExport, { called }] = useMutation<exportAllFaces>(EXPORT_ALL_FACES)

  return (
    <DeveloperToolsWrapper>
      <SectionTitle nospace>
        {t('settings.developer_tools.title', 'Under Construction')}
      </SectionTitle>
      <label htmlFor="dev_tools_export_all_faces_button">
        <InputLabelTitle>
          {t('settings.developer_tools.export_all.title', 'Export Portraits')}
        </InputLabelTitle>
        <InputLabelDescription>
          {t(
            'settings.developer_tools.export_all.description',
            'Export all faces to portrait folder'
          )}
        </InputLabelDescription>
      </label>
      <ExportAllFacesButton
        id="dev_tools_export_all_faces_button"
      />
    </DeveloperToolsWrapper>
  )
}

export default DeveloperTools
