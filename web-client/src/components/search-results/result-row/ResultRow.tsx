import * as React from 'react'
import {Component, ReactNode} from 'react'
import Highlight from "react-highlight";

import './ResultRow.scss'

export class ResultRow extends Component<{}, {}> {
    public render(): ReactNode {
        return (
            <div>
                <div>Path: /path/to/repository/file</div>
                <div>Last seen: 21 07 2019 - 11:33</div>
                <div>Commit: dc2105a</div>
                <Highlight>{`
var DEFAULT_LISTEN_ADDRESS = "127.0.0.1:7777"

func LoadConfig() (Config, error) {
                    nearestGitRepo, err := getNearestRepoPath()
                    if err != nil {
                    return *new(Config), err
}`}</Highlight>
            </div>
        );
    }
}
