#!/bin/bash

GIT_DIR=.git
GITHUB_DIR=.github
CUR_DIR=.
UP_DIR=..
function getdir(){
    for element in `ls -a $1`
    do  
        dir_or_file=$1"/"$element
        if [ -d $dir_or_file ]
        then
            if [[ "$element" == "$GIT_DIR" ]] || [[ "$element" == "$GITHUB_DIR" ]]
            then
                echo $dir_or_file
                rm -rf $dir_or_file
            else
                if [[ "$element" != "$CUR_DIR" ]] && [[ $element != "$UP_DIR" ]]
                then
                    getdir $dir_or_file
                fi    
            fi   
        fi  
    done
}
root_dir="../src/vendor"
getdir $root_dir