<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-alert type="success">工具类位置: src/utils/common-tools.ts</a-alert>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>小驼峰转小写下划线: toUnderline</div>
        <div>示例: userNameInfo -> user_name_info</div>
        <div>
          使用方式：
          <a-typography-text code>toUnderline(需要转换的字符串)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>下划线命名转小驼峰: getCamelCase</div>
        <div>示例: user_name_info -> userNameInfo</div>
        <div>
          使用方式：
          <a-typography-text code>getCamelCase(需要转换的字符串)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>生成随机Hex颜色: getColorHexColor</div>
        <div>该函数会随机生成并返回一个Hex的色值，如：#1eb31</div>
        <div>
          使用方式：
          <a-typography-text code>getColorHexColor()</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>判断变量数据类型: getObjType</div>
        <div>
          该函数会判断传入的变量类型并返回，如: boolean / number / string / function / object / array / date / regExp / undefined
          / null / object / symbol
        </div>
        <div>
          使用方式：
          <a-typography-text code>getObjType(需要判断的变量)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>基本数据类型数组去重: arrSet</div>
        <div>该函数会对基本数据类型数组去重，返回去重后的数组</div>
        <div>
          使用方式：
          <a-typography-text code>arrSet(需要去重的数组)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>数组对象去重: arrObjSet</div>
        <div>该函数会对数组对象去重，返回去重后的数组</div>
        <div>
          使用方式：
          <a-typography-text code>arrObjSet(需要去重的数组)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>数组对象根据指定key去重: somethingSet</div>
        <div>该函数会根据指定key对数组对象去重，返回去重后的数组</div>
        <div>
          使用方式：
          <a-typography-text code>somethingSet(需要去重的数组, 指定的key)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>排序-根据条件对数组对象排序: toSort</div>
        <div>该函数可根据指定key对数组对象排序，可设定升序、降序</div>
        <div>
          使用方式：
          <a-typography-text code>toSort(需要去重的数组, 指定key, 升序或降序)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>比对新对象和旧对象的差异值，返回新对象的差异值: Difference</div>
        <div>根据传入的新、旧Objct进行比对，返回他们的差异值(只比对顶层)</div>
        <div>
          使用方式：
          <a-typography-text code>Difference(新Object, 旧Object)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>删除数组中多个指定元素: arrMoreDeletion</div>
        <div>根据指定的key比对是否匹配，匹配则删除，最后返回删除后的数组，如：arrMoreDeletion(arr, [2,3,4], 'id')</div>
        <div>
          使用方式：
          <a-typography-text code>arrMoreDeletion(原数组, [key的对应值], 指定的key)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>无条件A/B数组交叉合并: crossList</div>
        <div>该函数会交叉合并A/B数组，返回合并后的数组</div>
        <div>
          使用方式：
          <a-typography-text code>crossList(A数组, B数组)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>生成指定长度的随机字符串: randomString</div>
        <div>包含数字、英文大小写，示例：randomString(5)</div>
        <div>
          使用方式：
          <a-typography-text code>randomString(指定随机字符串长度)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>获取数值在数组中的近似值: closest</div>
        <div>根据传入的值匹配数字数组中值最接近的数字并返回</div>
        <div>
          使用方式：
          <a-typography-text code>closest(数字数组, 需要匹配的值)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>根据标准值计算当前值的涨幅和减幅: getPercent</div>
        <div>
          该函数会根据传入的当前值、标准值进行计算，基于标准值计算出涨幅和减幅幅度，返回
          <a-typography-text code>{ percent: 10, type: 1, text: '涨幅10.00%'}</a-typography-text>
          格式
        </div>
        <div>
          使用方式：
          <a-typography-text code>getPercent(当前值, 标准值)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>根据涨幅和减幅计算实际值: getValue</div>
        <div>该函数会根据涨幅或减幅计算最终结果，示例：getValue(100, {type: 1, percent: 10})</div>
        <div>
          使用方式：
          <a-typography-text code>getValue(当前值, {type: 1涨幅2减幅, percent: 幅度百分比})</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>判断是否为空对象: isEmptyObject</div>
        <div>该函数会判断传入参数是否是空对象，返回true/false，区分Object/Array等类型</div>
        <div>
          使用方式：
          <a-typography-text code>isEmptyObject(传入的值)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>判断是否为安全环境: isSecureEnvironment</div>
        <div>返回true/false，安全环境为：https / localhost / 127.0.0.1 / 0.0.0.0</div>
        <div>
          使用方式：
          <a-typography-text code>isSecureEnvironment()</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>获取浏览器默认语言: webDefaultLanguage</div>
        <div>返回字符串，如：zh-CN</div>
        <div>
          使用方式：
          <a-typography-text code>webDefaultLanguage()</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>深拷贝变量: deepClone</div>
        <div>返回深拷贝后的值</div>
        <div>
          使用方式：
          <a-typography-text code>deepClone(需要深拷贝的变量)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>获取URL附带的参数: getQueryParams</div>
        <div>解析url?之后的参数转化为object的数据，默认取当前网页的url</div>
        <div>
          使用方式：
          <a-typography-text code>getQueryParams(url)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>截断字符串: truncateString</div>
        <div>截断字符串，如果字符串超过指定长度，则截断并添加省略号</div>
        <div>
          使用方式：
          <a-typography-text code>truncateString(字符串, 需要截断的长度)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>去掉字符串前后空格: verifyAndSpace</div>
        <div>该函数会去掉并返回字符串的前后空格</div>
        <div>
          使用方式：
          <a-typography-text code>verifyAndSpace(字符串)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>去掉字符串的所有中文和空格: verifyCnAndSpace</div>
        <div>该函数会去掉字符串的所有中文和空格并返回新字符串</div>
        <div>
          使用方式：
          <a-typography-text code>verifyCnAndSpace(字符串)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>去掉前后空格和字符串中的英文: verifyEnAndSpace</div>
        <div>该函数会返回纯非英文不包含前后空格的字符串</div>
        <div>
          使用方式：
          <a-typography-text code>verifyEnAndSpace(字符串)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>数字转中文大写: verifyNumberCnUppercase</div>
        <div>根据传入的数字转成中文大写字符串并返回，如：verifyNumberCnUppercase(1506045) -> 壹佰伍拾万陆仟零肆拾伍元整</div>
        <div>
          使用方式：
          <a-typography-text code>verifyNumberCnUppercase(数字)</a-typography-text>
        </div>
      </a-space>
    </div>
  </div>
</template>
