import { useMutation, useQuery, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'
import { TextField } from '../../primitives/form/Input'
import {
  InputLabelDescription,
  InputLabelTitle,
} from './SettingsPage'
import { classifyFaceThresholdQuery } from './__generated__/classifyFaceThresholdQuery'

export const CLASSIFY_FACE_THRESHOLD_QUERY = gql`
  query classifyFaceThresholdQuery {
    siteInfo {
      classifyFaceThreshold
    }
  }
`

export const ClassifyFaceThresholdInput = () => {
  const { t } = useTranslation()
  const threshold = useQuery<classifyFaceThresholdQuery>(CLASSIFY_FACE_THRESHOLD_QUERY) /* , {
    onCompleted(data) {
      const queryClassifyThreshold = data.siteInfo.classifyFaceThreshold
    },
  })
*/
  return (
          <TextField
            id="classify_face_Threshold"
            aria-label="Interval value"
            disabled={classifyFaceThresholdQuery.loading}
            value={classifyFaceThresholdQuery.value}
          />
  )
}

