# PaCo论文全文翻译摘要

- **卡片ID**: K024 | **版本**: v2
- **分类**: FHE核心技术模块 / Bootstrapping/PaCo
- **source_id**: `src_000024`
- **来源**: `01_源文件库/04_FHE核心技术模块/翻译.md`
- **置信度**: high | **需审核**: False
- **声明**: 翻译材料

## 摘要
PaCo(Partial Decryption-Based Bootstrapping)是CKKS bootstrapping的最新优化方案，核心创新：(1)结构化密钥——将密钥s拆分为h个Hamming weight=B的块s_j，每个s_j仅B=N/(4h)个非零位；(2)盲旋转解密——用test vector t_i(X)和旋转密钥同态解密每个s_j块；(3)多项式打包——将h个解密结果打包为多项式乘积。实现1.5×-2.6×加速。注意：PaCo的p=明文模数(plaintext modulus)，Δ=编码因子(scale factor)，两者不同。

## 关键知识点

### K024-1: PaCo密钥结构
- `s=s₁+…+s_h, 每个s_j有B=N/(4h)个非零位`
- 密钥拆分→每块稀疏→盲旋转只需B次旋转而非N次。h=2时B=N/8。

### K024-2: 盲旋转解密
- `t_i(X)=Σ_{j} test_j(X^{j})·Enc(解密密钥成分)`
- test vector预计算：包含所有可能解密值的加密结果。盲旋转选择匹配test vector的正确分量。

### K024-3: 多项式打包
- `Pack(m₁,…,m_h)=m₁·m₂·…·m_h`
- 将h个块解密结果打包为单一多项式，通过乘积Packing Relation连接。

### K024-4: p vs Δ 区别
- `p=明文模数(∼2^30), Δ=编码因子(∼2^80)`
- p是消息空间的模数（用于消息缩放），Δ是CKKS的编码放大因子。两者不是同一个量。

### K024-5: 实现性能
- `1.5×-2.6×加速 vs 之前的bootstrapping 方法`
- 在N=2^16, h=2, logQP≈1500时达到最佳加速比。密钥稀疏性直接贡献加速。

## 重要公式
- `s=Σ_{j=1}^{h}s_j, wt(s_j)=B=N/(4h)` — PaCo结构化密钥
- `t_i(X)=Σ_{e∈supp(s_j)}ct_{coeff_e}·X^e` — 盲旋转test vector

## 术语
- **盲旋转** (blind rotation): 用同态旋转和test vector实现的不暴露密钥的解密
- **多项式打包** (polynomial packing): 将多个解密块结果合并为一个多项式
- **结构化密钥** (structured secret key): 将密钥拆分为稀疏块的技巧

## 关系
- abstract_algebra: -
- RLWE_lattice: -
- FHE_technique: CKKS Bootstrapping最新进展
- current_paper: 当前论文核心——Eurocrypt 2027论文的implementation部分

## RAG建议: 最高优先级，当前论文核心参考论文
