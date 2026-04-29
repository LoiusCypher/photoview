import { useMutation, gql } from '@apollo/client'
import React, { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '../../primitives/form/Input'
import { exportFaces } from './__generated__/exportFaces'
import { exportConfirmedFaces, exportConfirmedFacesVariables } from './__generated__/exportConfirmedFaces'

const EXPORT_ALL_FACES = gql`
  mutation exportFaces {
    exportAllFaces {
      success
      message
    }
  }
`

const EXPORT_FACES = gql`
  mutation exportConfirmedFaces($onlyConfirmed: Boolean!) {
    exportFaces(onlyConfirmed: $onlyConfirmed) {
      success
      message
    }
  }
`

export const ExportAllFacesButton = () => {
  const { t } = useTranslation()
  const [startExport, { called }] = useMutation<exportFaces>(EXPORT_ALL_FACES)
  const [startConfirmedExport, { calledExport }] = useMutation<exportConfirmedFaces, exportConfirmedFacesVariables>(EXPORT_FACES)

  return (
    <>
    <Button
      className="mb-4"
      onClick={() => {
        startExport()
      }}
      disabled={called}
    >
      {t('settings.export_all_faces', 'Export All Faces')}
    </Button>
    <Button
      className="mb-4"
      onClick={() => {
        startConfirmedExport({variables: {onlyConfirmed: true}})
      }}
      disabled={calledExport}
    >
      {t('settings.export_confirmed_faces', 'Export Only Confirmed Faces')}
    </Button>
    </>
  )
}

