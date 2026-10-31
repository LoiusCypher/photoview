import React from 'react'
import { authToken } from '../../../helpers/authentication'
import { useMutation, gql } from '@apollo/client'
import { useTranslation } from 'react-i18next'
import { Button } from '../../../primitives/form/Input'
import classNames from 'classnames'
import { orderFaceGroupAction, orderFaceGroupVariables } from './__generated__/orderFaceGroupAction'
import { searchFaceGroupAction, searchFaceGroupVariables } from './__generated__/searchFaceGroupAction'
import { splitFaceGroupAction, splitFaceGroupVariables } from './__generated__/splitFaceGroupAction'

const ORDER_FACE_GROUP_MUTATION = gql`
  mutation orderFaceGroupAction( $groupId: ID!) {
    orderFacesOfFaceGroup( groupId: $groupId) {
      success
      message
    }
  }
`

const SEARCH_FACE_GROUP_MUTATION = gql`
  mutation searchFaceGroupAction( $groupId: ID!) {
    searchForFaceGroup( groupId: $groupId) {
      success
      message
    }
  }
`

const SPLIT_FACE_GROUP_MUTATION = gql`
  mutation splitFaceGroupAction( $groupId: ID!) {
    splitFaceGroup( groupId: $groupId) {
      success
      message
    }
  }
`

export type FaceGroupOrderButtonProps = {
  groupID: int,
}

export type FaceGroupSearchButtonProps = {
  groupID: int,
}

export type FaceGroupSplitButtonProps = {
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
      disabled={calledOrder}
    >
      {t('face_group.order', 'ReOrder / Exec MoveTo face group')}
    </Button>
  )
}

const FaceGroupSearchButton = ({
  groupID,
}: FaceGroupSearchButtonProps) => {
  const { t } = useTranslation()
  const [startFaceGroupSearch, { calledSearch }] = useMutation<searchFaceGroupAction,searchFaceGroupVariables>(SEARCH_FACE_GROUP_MUTATION)

  return (
    <Button
      className="mb-1"
      onClick={() => { startFaceGroupSearch({variables: {groupId: groupID}});}}
      disabled={calledSearch}
    >
      {t('face_group.search', 'Search others / Check MoveTo face group')}
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
      {t('face_group.split', 'Split group')}
    </Button>
  )
}

type FaceGroupOrderProps = {
  groupID: int,
}

type FaceGroupSearchProps = {
  groupID: int,
}

type FaceGroupSplitProps = {
  groupID: int,
}

const FaceGroupRow = ({
  groupID,
}: FaceGroupSplitProps) => {
  return (
    <div className="flex items-end gap-4 flex-wrap mb-4">
      {authToken() && (
        <FaceGroupOrderButton
          groupID={groupID}
        />
      )}
      {authToken() && (
        <FaceGroupSearchButton
          groupID={groupID}
        />
      )}
      {authToken() && (
        <FaceGroupSplitButton
          groupID={groupID}
        />
      )}
    </div>
  )
}

export default FaceGroupRow
