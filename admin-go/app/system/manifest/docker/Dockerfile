FROM loads/alpine:3.8

###############################################################################
#                                INSTALLATION
###############################################################################

ENV WORKDIR                 /app
ADD resource                $WORKDIR/
ADD ./temp/linux_amd64/main $WORKDIR/main
ADD ./temp/linux_amd64/migrate $WORKDIR/migrate
ADD ./temp/linux_amd64/database $WORKDIR/database
ADD ./manifest/docker/start.sh $WORKDIR/start.sh
RUN chmod +x $WORKDIR/main $WORKDIR/migrate $WORKDIR/start.sh

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./start.sh
