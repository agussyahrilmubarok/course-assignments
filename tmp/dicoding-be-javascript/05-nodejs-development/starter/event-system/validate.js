const { execSync } = require('child_process');
const { join } = require('path');
const fs = require('fs');

(() => {
  try {

    execSync('node -c  app.js', { cwd: __dirname, stdio: 'pipe' });

    const stdoutAppJs = execSync('node app.js', { cwd: __dirname, encoding: 'utf-8', stdio: 'pipe' });

    if (stdoutAppJs.toString() !== '5\n-1\n11\n') {
      throw new Error('Berkas `app.js` tidak menghasilkan output yang sesuai. Harap periksa kembali, ya!')
    }

    try {
      const additionalCodeToAdd = 'CmNvbnN0IGFzc2VydCA9IHJlcXVpcmUoJ2Fzc2VydCcpLnN0cmljdDsKCmZ1bmN0aW9uIGlzT2JqZWN0TGlrZUV2ZW50RW1pdHRlcihvYmplY3QpIHsKICByZXR1cm4gdHlwZW9mIG9iamVjdC5lbWl0ID09PSAnZnVuY3Rpb24nCiAgICAmJiB0eXBlb2Ygb2JqZWN0Lm9uID09PSAnZnVuY3Rpb24nCiAgICAmJiB0eXBlb2Ygb2JqZWN0Lm9uY2UgPT09ICdmdW5jdGlvbicKICAgICYmIHR5cGVvZiBvYmplY3Qub2ZmID09PSAnZnVuY3Rpb24nOwp9Cgp0cnkgewogIGFzc2VydC5zdHJpY3RFcXVhbCh0eXBlb2YgQmFzaWNNYXRoLCAnb2JqZWN0JywgJ1Bhc3Rpa2FuIGthbXUgdGlkYWsgbWVuZ3ViYWggYXBhIHB1biBwYWRhIGtvZGUgYEJhc2ljTWF0aGAuIEhhcmFwIHVudHVrIGRpY2VrIGxhZ2ksIHlhISBeX14nKTsKICBhc3NlcnQuc3RyaWN0RXF1YWwodHlwZW9mIEV2ZW50RW1pdHRlciwgJ2Z1bmN0aW9uJywgJ1Bhc3Rpa2FuIGthbXUgbWVuZ2ltcG9yIGV2ZW50IGVtaXR0ZXIgZGVuZ2FuIG5hbWEgdmFyaWFiZWwgYEV2ZW50RW1pdHRlcmAuJyk7CiAgYXNzZXJ0LnN0cmljdEVxdWFsKGV2ZW50RW1pdHRlciBpbnN0YW5jZW9mIEV2ZW50RW1pdHRlciwgdHJ1ZSwgJ1Bhc3Rpa2FuIGthbXUgbWVtYnVhdCBpbnN0YW5jZSBkYXJpIGBFdmVudEVtaXR0ZXJgIGRlbmdhbiBuYW1hIGBldmVudEVtaXR0ZXJgLiBKYW5nYW4gc2FtcGFpIHNhbGFoLCB5YSEgXl9eJyk7CiAgYXNzZXJ0LnN0cmljdEVxdWFsKGlzT2JqZWN0TGlrZUV2ZW50RW1pdHRlcihldmVudEVtaXR0ZXIpLCB0cnVlLCAnUGFzdGlrYW4gYGV2ZW50RW1pdHRlcmAgYWRhbGFoIGluc3RhbmNlIGRhcmkgYEV2ZW50RW1pdHRlcmAuIEphbmdhbiBzYW1wYWkga2VsaXJ1LCB5YSEgXl9eJyk7CiAgYXNzZXJ0LnN0cmljdEVxdWFsKGV2ZW50RW1pdHRlci5saXN0ZW5lcnMoJ2NhbGN1bGF0ZScpLmxlbmd0aCwgMSwgJ1Bhc3Rpa2FuIGthbXUgbWVuZGFmdGFya2FuIGBCYXNpY01hdGguYWRkYCBzZWJhZ2FpIGxpc3RlbmVyIGRhcmkgZXZlbnQgImNhbGN1bGF0ZSIuIEphbmdhbiBzYW1wYWkgbHVwYSwgeWEhIF5fXicpOwogIGFzc2VydC5zdHJpY3RFcXVhbChldmVudEVtaXR0ZXIubGlzdGVuZXJzKCdjYWxjdWxhdGUnKVswXSA9PT0gQmFzaWNNYXRoLmFkZCwgdHJ1ZSwgJ1Bhc3Rpa2FuIGthbXUgbWVuZGFmdGFya2FuIGBCYXNpY01hdGguYWRkYCBzZWJhZ2FpIGxpc3RlbmVyIGRhcmkgZXZlbnQgImNhbGN1bGF0ZSIuIEphbmdhbiBzYW1wYWkgbHVwYSwgeWEhIF5fXicpOwp9IGNhdGNoIChlcnJvcikgewogIGlmIChlcnJvci5uYW1lID09PSAnQXNzZXJ0aW9uRXJyb3InKSB7CiAgICBjb25zb2xlLmxvZyhlcnJvci5tZXNzYWdlKTsKICAgIHByb2Nlc3MuZXhpdCgxKTsKICB9IGVsc2UgewogICAgY29uc29sZS5sb2coYEthbWkgZ2FnYWwgbWVuZ3VqaSBrb2RlIHlhbmcga2FtdSB0dWxpcy4gQmVyaWt1dCBlcnJvciB5YW5nIGRpaGFzaWxrYW46ICR7ZXJyb3IubWVzc2FnZX0gU2ltYWsgZXJyb3IgdGVyc2VidXQgZGFuIHBhc3Rpa2FuIGthbXUgc3VkYWggbWVudWxpc2thbiBrb2RlIGRlbmdhbiBiZW5hciwgeWEhYCk7CiAgICBwcm9jZXNzLmV4aXQoMSk7CiAgfQp9';
      const appJSCurrentCode = fs.readFileSync(join(__dirname, 'app.js'), 'utf-8');

      fs.writeFileSync(join(__dirname, 'app.testing.js'), appJSCurrentCode + Buffer.from(additionalCodeToAdd, 'base64').toString('utf-8'));

      execSync('node app.testing.js', { cwd: __dirname, stdio: 'pipe' });
      console.log('passed!')
    } catch (error) {
      const { stderr } = error;
      console.log(stderr.toString());
    } finally {
      fs.unlinkSync(join(__dirname, 'app.testing.js'));
    }
  } catch (error) {
    console.log(error.message);
  }
})()
