const fs = require('fs')
const path = require('path')

let opts = {
    ext: ['go'],
    ignore: ['vendor',
        'example',
        'tests',
        '.git',
        '.idea',
        'go.mod',
        'go.sum',
        'config.example.json',
        'example.json',
        'notifier.proto',
        'notifier.pb.go',
        'notifier_grpc.pb.go',
        'purpose.txt',
        'stats.txt'
    ],
    path: '../',
    commentsPrefix: ['//']
}

function countLines(opts) {
    if(!opts.result)
        opts.result = {
            lines: 0,
            comments: 0
        }

    let stats = fs.statSync(opts.path)

    if(stats.isDirectory()) {
        let files = fs.readdirSync(opts.path)
        for(let file of files) {
            if(opts.ignore.includes(file.trim())) {
                continue
            }
            let _path = path.join(opts.path, file)
            console.log('reading:', _path)
            let newOpts = {
                ext: opts.ext,
                ignore: opts.ignore,
                path: _path,
                result: opts.result,
                commentsPrefix: opts.commentsPrefix
            }
            countLines(newOpts);
        }
    } else {
        
        let stream = fs.readFileSync(opts.path)

        let lines = stream.toString().split('\n')

        for(let line of lines) {
            let normalLine = line.trim();

            let isComment = false;
            for(let prefix of opts.commentsPrefix) {
                if(normalLine.startsWith(prefix)) {
                    isComment = true;
                }
                break;
            }
            if(isComment) {
                opts.result.comments++;
                
            } else {
                opts.result.lines++
            }
        }
    }

    return opts.result
}

console.log(countLines(opts))