import { useMutation, gql } from '@apollo/client'
import React, { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '../../primitives/form/Input'
import { exportFaces } from './__generated__/exportFaces'

const EXPORT_ALL_FACES = gql`
  mutation exportFaces {
    exportAllFaces {
      success
      message
    }
  }
`

export const ExportAllFacesButton = () => {
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

