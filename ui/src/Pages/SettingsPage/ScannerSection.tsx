import React from 'react'
import PeriodicScanner from './PeriodicScanner'
import { ScannerConcurrentWorkers } from './ScannerConcurrentWorkers'
import { ScannerClassifyThreshold } from './ScannerClassifyThreshold'
import { ScannerScanAllUsers } from './ScannerScanAllUsers'
import { SectionTitle } from './SettingsPage'
import { useTranslation } from 'react-i18next'

const ScannerSection = () => {
  const { t } = useTranslation()

  return (
    <div>
      <SectionTitle nospace>
        {t('settings.scanner.title', 'Scanner')}
      </SectionTitle>
      <ScannerClassifyThreshold />
      <ScannerScanAllUsers />
      <PeriodicScanner />
      <ScannerConcurrentWorkers />
    </div>
  )
}

export default ScannerSection
