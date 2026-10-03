import { useMutation, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'
import { Button } from '../../primitives/form/Input'
import {
  InputLabelDescription,
  InputLabelTitle,
} from './SettingsPage'
import { checkGroups, checkGroupsVariables } from './__generated__/checkGroups'

const CHECK_FACE_GROUP = gql`
  mutation checkGroups($faceGroupID: ID!) {
    checkFaceGroup(faceGroupID: $faceGroupID) {
      success
      message
    }
  }
`

export const CheckFaceGroupButton = () => {
  const { t } = useTranslation()
  const [startCheck, { called }] = useMutation<checkGroups, checkGroupsVariables>(CHECK_FACE_GROUP, { variables: { faceGroupID: '2', }, } )

  return (
    <div>
      <Button
        className="mb-4"
        onClick={() => {
          startCheck()
        }}
        disabled={called}
      >
        {t('settings.check_face_group', 'Check Face Group')}
      </Button>
      { called && (
        <>
        <span className="sr-only">Loading...</span>
      <svg aria-hidden="true" className="inline-block w-8 h-8 text-neutral-tertiary animate-rotate" viewBox="-2 0 100 101" xmlns="http://www.w3.org/2000/svg">
        <circle className="animate-dash" cx="50" cy="50" r="45">
        </circle>
      </svg>
      </>
      )}
    </div>
  )
}

