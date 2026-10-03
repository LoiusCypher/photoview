import React from 'react'
import { authToken } from '../../helpers/authentication'
import { useMutation, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'
import { Button } from '../../primitives/form/Input'
import classNames from 'classnames'
import { scanAlbumAction, scanAlbumVariables } from './__generated__/scanAlbumAction'

const SCAN_ALBUM_MUTATION = gql`
  mutation scanAlbumAction( $albumId: ID!) {
    scanAlbum( albumId: $albumId) {
      success
      message
    }
  }
`

export type RescanAlbumButtonProps = {
  albumId: int,
}

export const RescanAlbumButton = ({
  albumId,
}: RescanAlbumButtonProps) => {
  const { t } = useTranslation()
  const [startAlbumScanner, { calledAlbum }] = useMutation<scanAlbumAction,scanAlbumVariables>(SCAN_ALBUM_MUTATION)

  return (
    <Button
      className="mb-1"
      onClick={() => { startAlbumScanner({variables: {albumId: albumId}});}}
      disabled={calledAlbum}
    >
      {t('album_rescan.only_favorites', 'Rescan album for faces')}
    </Button>
  )
}

type AlbumRescanProps = {
  albumId: int,
}

const AlbumRescan = ({
  albumId,
}: AlbumRescanProps) => {
  return (
    <div className="flex items-end gap-4 flex-wrap mb-4">
      {authToken() && (
        <RescanAlbumButton
          albumId={albumId}
        />
      )}
    </div>
  )
}

export default AlbumRescan
