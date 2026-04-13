import React from 'react'
import { useMutation, gql } from '@apollo/client'
import PeriodicScanner from './PeriodicScanner'
import { ScannerConcurrentWorkers } from './ScannerConcurrentWorkers'
import { ScannerClassifyThreshold } from './ScannerClassifyThreshold'
import { SectionTitle, InputLabelDescription } from './SettingsPage'
import { useTranslation } from 'react-i18next'
import { scanAllMutation } from './__generated__/scanAllMutation'
import { Button } from '../../primitives/form/Input'

import { TextField } from '../../primitives/form/Input'

const SCAN_MUTATION = gql`
  mutation scanAllMutation {
    scanAll {
      finished
      success
      message
    }
  }
`

const ScannerSection = () => {
  const { t } = useTranslation()
  const [startScanner, { data: data, loading, error, called }] = useMutation<scanAllMutation>(SCAN_MUTATION)

  console.log('apollo', data?.scanAll);

  return (
    <div>
      <SectionTitle nospace>
        {t('settings.scanner.title', 'Scanner')}
      </SectionTitle>
      <ScannerClassifyThreshold />
      <InputLabelDescription>
        {t(
          'settings.scanner.description',
          'Will scan all users for new or updated media'
        )}
      </InputLabelDescription>
      <div className="flex">
      <div className="relative">
        <Button
          onClick={() => {
            startScanner()
          }}
          disabled={(loading || called) && !data}
        >
          {t('settings.scanner.scan_all_users', 'Scan all users')}
        </Button>
        { (loading || called) && !data && (
        <>
          <span className="sr-only">Loading...</span>
          <svg aria-hidden="true" className="absolute top-0 left-0 w-full h-full text-neutral-tertiary animate-rotate" viewBox="-2 0 100 101" xmlns="http://www.w3.org/2000/svg">
            <circle className="animate-dash" cx="50" cy="50" r="45">
            </circle>
          </svg>
        </>
        )}
      </div>
      <p className="inline-block">
      </p>
      </div>
      <PeriodicScanner />
      <ScannerConcurrentWorkers />
    </div>
  )
}

export default ScannerSection
