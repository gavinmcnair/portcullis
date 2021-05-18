FROM golang:1.15.6 as builder

WORKDIR .

COPY . .

RUN make portcullis

FROM scratch
MAINTAINER Gavin McNair

ARG git_repository="Unknown"
ARG git_commit="Unknown"
ARG git_branch="Unknown"
ARG built_on="Unknown"

LABEL git.repository=$git_repository
LABEL git.commit=$git_commit
LABEL git.branch=$git_branch
LABEL build.on=$built_on

COPY --from=builder /portcullis/bin/linux/portcullis .

CMD [ "/portcullis" ]
