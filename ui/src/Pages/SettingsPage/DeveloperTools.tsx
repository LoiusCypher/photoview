import { useMutation, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'
import { Button } from '../../primitives/form/Input'
import {
  InputLabelDescription,
  InputLabelTitle,
  SectionTitle,
} from './SettingsPage'
import { checkGroups, checkGroupsVariables } from './__generated__/checkGroups'
import { classifyFaceThresholdQuery } from './__generated__/classifyFaceThresholdQuery'
import { ExportAllFacesButton } from './Dev_ExportFaces'

const CHECK_FACE_GROUP = gql`
  mutation checkGroups($faceGroupID: ID!) {
    checkFaceGroup(faceGroupID: $faceGroupID) {
      success
      message
    }
  }
`

export const CLASSIFY_FACE_THRESHOLD_QUERY = gql`
  query classifyFaceThresholdQuery {
    siteInfo {
      classifyFaceThreshold
    }
  }
`

const CheckFaceGroupButton = () => {
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
      <label htmlFor="dev_tools_check_face_gruop_button">
        <InputLabelTitle>
          {t('settings.developer_tools.check_group.title', 'Check Face Group')}
        </InputLabelTitle>
        <InputLabelDescription>
          {t(
            'settings.developer_tools.check_group.description',
            'Check face goup member consistancy'
          )}
        </InputLabelDescription>
      </label>
      <CheckFaceGroupButton
        id="dev_tools_check_face_gruop_button"
      />
    </DeveloperToolsWrapper>
  )
}

export default DeveloperTools
