import { useMutation, useQuery } from '@apollo/client'
import React, { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'
import {
  SectionTitle,
} from './SettingsPage'
import { exportAllFaces } from './__generated__/exportFaces'

const EXPORT_ALL_FACES = gql`
  mutation exportFaces {
    exportAllFaces {
      success
      message
    }
  }
`

const DeveloperTools = () => {
  const { t } = useTranslation()

  return (
    <div>
      <SectionTitle nospace>
        {t('settings.user_preferences.title', 'User preferences')}
      </SectionTitle>
    </div>
  )
}

export default DeveloperTools
