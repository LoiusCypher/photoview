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
    <Button
      className="mb-4"
      onClick={() => {
        startCheck()
      }}
      disabled={called}
    >
      {t('settings.check_face_group', 'Check Face Group')}
    </Button>
  )
}

