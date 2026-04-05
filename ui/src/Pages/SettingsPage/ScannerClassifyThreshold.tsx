import React, { useRef, useState } from 'react'
import { useQuery, useMutation, gql } from '@apollo/client'
import { InputLabelTitle, InputLabelDescription } from './SettingsPage'
import { useTranslation } from 'react-i18next'
import { classifyFaceThresholdQuery } from './__generated__/classifyFaceThresholdQuery'
import {
  setClassifyThresholdMutation,
  setClassifyThresholdMutationVariables,
} from './__generated__/setClassifyThresholdMutation'
import { TextField } from '../../primitives/form/Input'
import './style-spin.css'

export const CLASSIFY_FACE_THRESHOLD_QUERY = gql`
  query classifyFaceThresholdQuery {
    siteInfo {
      classifyFaceThreshold
    }
  }
`

export const SET_CLASSIFY_THRESHOLD_MUTATION = gql`
  mutation setClassifyThresholdMutation($threshold: Float!) {
    setFaceClassifyThreshold(threshold: $threshold)
  }
`

export const ScannerClassifyThreshold = () => {
  const { t } = useTranslation()

  const [threshold, setThreshold] = useState(-1)
  const classifyThresholdQuery = useQuery<classifyFaceThresholdQuery>(CLASSIFY_FACE_THRESHOLD_QUERY,
   {
     onCompleted(data) {
       setThreshold(data.siteInfo.classifyFaceThreshold)
     },
   }
)

  const [setThresholdMutation, setThresholdMutationData] = useMutation<
    setClassifyThresholdMutation,
    setClassifyThresholdMutationVariables
  >(SET_CLASSIFY_THRESHOLD_MUTATION, { refetchQueries: [ CLASSIFY_FACE_THRESHOLD_QUERY ], })

  const updateClassifyThreshold = (classifyThreshold: number) => {
    setThresholdMutation({
      variables: {
        threshold: classifyThreshold,
      },
    })
  }

  if (classifyThresholdQuery.isError || !classifyThresholdQuery.data) {
    return <span>Error. Please reload page.</span>;
  }

  return (
    <div>
      <label htmlFor="scanner_classify_threshold_field">
        <InputLabelTitle>
          {t('settings.classify_faces_threshold.title', 'Classify faces threshold')}
        </InputLabelTitle>
        <InputLabelDescription>
          {t(
            'settings.classify_faces_threshold.description',
            'The maximum distance to consider faces to belong to the same person'
          )}
        </InputLabelDescription>
      </label>
      { (classifyThresholdQuery.loading || setThresholdMutationData.loading) &&
      <div>
        <span className="sr-only">Loading...</span>
        <svg aria-hidden="true" className="w-8 h-8 text-neutral-tertiary animate-rotate" viewBox="-2 0 100 101" xmlns="http://www.w3.org/2000/svg">
          <circle className="animate-dash" cx="50" cy="50" r="45">
          </circle>
        </svg>
      </div>
      }
      { (classifyThresholdQuery.loading || setThresholdMutationData.loading) ||
      <TextField
        disabled={classifyThresholdQuery.loading || setThresholdMutationData.loading}
        type="number"
        min="-1"
        max="0.5"
	step="0.01"
        id="scanner_classify_threshold_field"
        value={threshold}
        onChange={event => {
          if (!isNaN(parseFloat(event.target.value))) {
            var thresh = parseFloat(event.target.value)
            setThreshold(thresh);
            if (event.target !== document.activeElement) {
              updateClassifyThreshold(thresh)
            }
          }
        }}
        onBlur={() => { updateClassifyThreshold(threshold);}}
        onKeyDown={event =>
          event.key == 'Enter' && updateClassifyThreshold(this.value)
        }
      />
      }
    </div>
  )
}
