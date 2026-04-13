import React from 'react'
import { authToken } from '../../../helpers/authentication'
import { useMutation, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'
import { Button } from '../../../primitives/form/Input'
import classNames from 'classnames'
import { splitFaceGroupAction, splitFaceGroupVariables } from './__generated__/splitFaceGroupAction'

const SPLIT_FACE_GROUP_MUTATION = gql`
  mutation splitFaceGroupAction( $groupId: ID!) {
    splitFaceGroup( groupId: $groupId) {
      success
      message
    }
  }
`

export type FaceGroupSplitButtonProps = {
  groupID: int,
}

const FaceGroupSplitButton = ({
  groupID,
}: FaceGroupSplitButtonProps) => {
  const { t } = useTranslation()
  const [startFaceGroupSplit, { calledSplit }] = useMutation<splitFaceGroupAction,splitFaceGroupVariables>(SPLIT_FACE_GROUP_MUTATION)

  return (
    <Button
      className="mb-1"
      onClick={() => { startFaceGroupSplit({variables: {groupId: groupID}});}}
      disabled={calledSplit}
    >
      {t('face_group.spllit', 'Split face group')}
    </Button>
  )
}

type FaceGroupSplitProps = {
  groupID: int,
}

const FaceGroupSplit = ({
  groupID,
}: FaceGroupSplitProps) => {
  return (
    <div className="flex items-end gap-4 flex-wrap mb-4">
      {authToken() && (
        <FaceGroupSplitButton
          groupID={groupID}
        />
      )}
    </div>
  )
}

export default FaceGroupSplit
