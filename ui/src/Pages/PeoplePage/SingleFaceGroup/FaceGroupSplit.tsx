import React from 'react'
import { authToken } from '../../../helpers/authentication'
import { useMutation, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'
import { Button } from '../../../primitives/form/Input'
import classNames from 'classnames'
import { splitFaceGroupAction, splitFaceGroupVariables } from './__generated__/splitFaceGroupAction'
import { orderFaceGroupAction, orderFaceGroupVariables } from './__generated__/orderFaceGroupAction'

const SPLIT_FACE_GROUP_MUTATION = gql`
  mutation splitFaceGroupAction( $groupId: ID!) {
    splitFaceGroup( groupId: $groupId) {
      success
      message
    }
  }
`

const ORDER_FACE_GROUP_MUTATION = gql`
  mutation orderFaceGroupAction( $groupId: ID!) {
    orderFaceGroup( groupId: $groupId) {
      success
      message
    }
  }
`

export type FaceGroupSplitButtonProps = {
  groupID: int,
}

export type FaceGroupOrderButtonProps = {
  groupID: int,
}

const FaceGroupOrderButton = ({
  groupID,
}: FaceGroupOrderButtonProps) => {
  const { t } = useTranslation()
  const [startFaceGroupOrder, { calledOrder }] = useMutation<orderFaceGroupAction,orderFaceGroupVariables>(ORDER_FACE_GROUP_MUTATION)

  return (
    <Button
      className="mb-1"
      onClick={() => { startFaceGroupOrder({variables: {groupId: groupID}});}}
      disabled={calledSplit}
    >
      {t('face_group.order', 'Order face group')}
    </Button>
  )
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

type FaceGroupOrderProps = {
  groupID: int,
}

type FaceGroupSplitProps = {
  groupID: int,
}

const FaceGroupOrder = ({
  groupID,
}: FaceGroupOrderProps) => {
  return (
    <div className="flex items-end gap-4 flex-wrap mb-4">
      {authToken() && (
        <FaceGroupOrderButton
          groupID={groupID}
        />
      )}
    </div>
  )
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

export default FaceGroupSplit, FaceGroupOrder
