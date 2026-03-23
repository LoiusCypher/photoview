import React, { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQuery, useMutation, gql } from '@apollo/client'
import Checkbox from '../../primitives/form/Checkbox'
import { scanFacesOnOriginalFilesQuery } from './__generated__/scanFacesOnOriginalFilesQuery'

const SCAN_ON_ORGINAL_FILES_QUERY = gql`
  query scanFacesOnOriginalFilesQuery {
    siteInfo {
      scanFacesOnOriginalFiles
    }
  }
`

const ScanFacesOnOriginalFilesCheckbox = () => {
  const { t } = useTranslation()
  const { scanOriginalFiles } = useQuery<scanFacesOnOriginalFilesQuery>(SCAN_ON_ORGINAL_FILES_QUERY)

  return (
      <Checkbox
        label={t(
          'settings.periodic_scanner.checkbox_label',
          'Scan faces fromoriginal Files instead of thumbnails'
        )}
        disabled={scanFacesOnOriginalFilesQuery.loading}
        checked={scanOriginalFiles}
      />
  )
}

export default ScanFacesOnOriginalFilesCheckbox
