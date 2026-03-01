import { useMutation, useQuery, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'
import styled from 'styled-components'
import {
  InputLabelDescription,
  InputLabelTitle,
  SectionTitle,
} from './SettingsPage'
import { CheckFaceGroupButton } from './Dev_CheckFaceGroups'
import { ExportAllFacesButton } from './Dev_ExportFaces'
import ScanFacesOnOriginalFilesCheckbox from './Dev_ScanOriginalFiles'

const DeveloperToolsWrapper = styled.div`
  margin-bottom: 24px;
`

const DeveloperTools = () => {
  const { t } = useTranslation()

  return (
    <DeveloperToolsWrapper>
      <ScanFacesOnOriginalFilesCheckbox />
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
