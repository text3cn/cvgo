"use strict";
const infoConsole = document.getElementById('consoleInfo');
if (console) {
    const _console = {
        log: console.log,
    };
    console.log = function (attr) {
        _console.log(attr);
        const rewrite = function (obj) {
            for (const key in obj) {
                if (typeof obj[key] === 'object') {
                    rewrite(obj[key]);
                    continue;
                }
                obj[key] = `t3v_left${obj[key]}t3v_right`;
            }
        };
        rewrite(attr);
        const str = JSON.stringify(attr, null, 4);
        const node = document.createElement('H2');
        const textnode = document.createTextNode(str);
        node.appendChild(textnode);
        infoConsole.appendChild(node);
        let html = $('H2').html();
        html = html.replace(/".*":/gm, str => {
            let attrName = str.replaceAll('"', '');
            attrName = attrName.replace(':', '');
            return `<span class="attrName">${attrName} :</span>`;
        });
        html = html.replace(/t3v_left.*t3v_right/mg, str => {
            const ret = `<span class="attrValue">${str}</span>`;
            return ret;
        });
        html = html.replaceAll('t3v_left', '');
        html = html.replaceAll('t3v_right', '');
        html = html.replaceAll(':', '<span class="colon">:</span>');
        $('H2').html(html);
    };
}
//# sourceMappingURL=common.js.map