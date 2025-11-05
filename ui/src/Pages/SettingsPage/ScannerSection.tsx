import React from 'react'
import { useMutation, gql } from '@apollo/client'
import PeriodicScanner from './PeriodicScanner'
import { ScannerConcurrentWorkers } from './ScannerConcurrentWorkers'
import { SectionTitle, InputLabelDescription } from './SettingsPage'
import { useTranslation } from 'react-i18next'
import { scanAllMutation } from './__generated__/scanAllMutation'
import { watchModifiedQuery } from './__generated__/watchModifiedQuery'
import { Button } from '../../primitives/form/Input'
import Checkbox from '../../primitives/form/Checkbox'

const SCAN_MUTATION = gql`
  mutation scanAllMutation {
    scanAll {
      success
      message
    }
  }
`

export const WATCH_MODIFIED_QUERY = gql`
  query watchModifiedQuery {
    siteInfo {
      watchModifiedTime
    }
  }
`

const WatchModificationTime = () => {
  const { t } = useTranslation()

  const [enableWatchModificationTime, setEnableWatchModificationTime] = useState(false)

  const watchModificationTimeQuery = useQuery<watchModifiedQuery>(SCAN_INTERVAL_QUERY, {
    onCompleted(data) {
      const queryWatchModifiedQuery= data.siteInfo.watchModifiedTime

      setEnableWatchModificationTime(queryWatchModifiedQuery)
    },
  })
}

const ScannerSection = () => {
  const { t } = useTranslation()
  const [startScanner, { called }] = useMutation<scanAllMutation>(SCAN_MUTATION)

  const onWatchModifionTimeCheckboxChange = (checked: boolean) => {
    setEnableWatchModificationTime(checked)
  }

  return (
    <div>
      <SectionTitle nospace>
        {t('settings.scanner.title', 'Scanner')}
      </SectionTitle>
      <InputLabelDescription>
        {t(
          'settings.scanner.description',
          'Will scan all users for new or updated media'
        )}
      </InputLabelDescription>
      <Button
        onClick={() => {
          startScanner()
        }}
        disabled={called}
      >
        {t('settings.scanner.scan_all_users', 'Scan all users')}
      </Button>

      <h3 className="font-semibold text-lg mt-4 mb-2">
        {t('settings.scanner.watch_modified_times', 'Scanner watches media files modification time')}
      </h3>

      <Checkbox
        label={t(
          'settings.scanner.watch_modified_times.checkbox_label',
          'Watch modified times'
        )}
        //disabled={scanIntervalQuery.loading}
        checked={enableWatchModificationTime}
        onChange={event =>
          onWatchModifionTimeCheckboxChange(event.target.checked || false)
        }
      />
      <PeriodicScanner />
      <ScannerConcurrentWorkers />
    </div>
  )
}

export default ScannerSection
