# Batch-05a3 处理日志

> 2026-05-04T15:40 | SIMD / slot permutation / NTT 整理

## 范围
- ✅ SIMD编码 (BGV/BFV CRT-based)
- ✅ Slot编号的群论基础 (Frobenius轨道、Z_m^*/<p>)
- ✅ Slot Permutation机制 (tau_j自同构诱导陪集置换)
- ✅ tau_j: a(X)->a(X^j)自同构定义与槽视角
- ✅ NTT/DFT/FFT定义与区别
- ✅ Vandermonde/DFT/典范嵌入矩阵关系
- ✅ Psi_n,i: K_n -> L_i嵌入与第一同构定理
- ❌ 噪声/安全/参数
- ❌ 实验结果
- ❌ proof cleanup

## 新增导入 (14 files)

| source_id | 文件 | 标注 |
|-----------|------|------|
| src_000081 | 1slotpermu.md | slot编号+permutation |
| src_000082 | 2不同槽之间的本原m次根.md | Frobenius等价类 |
| src_000083 | t维度旋转.md | 多维槽索引 |
| src_000084 | 自同构.md | tau_j定义 |
| src_000085 | 自己的话.md | 简短总结 |
| src_000086 | Freobinus.md | d=ord_m(p) |
| src_000087 | SIMD操作.md | SIMD编码主文档 |
| src_000088 | 编码方法.md | Psi_n,i嵌入+Gamma_n,l |
| src_000089 | 代数：幂集与方程的解.md | psi根->基 |
| src_000090 | 代数：扩张与代值映射.md | 第一同构定理 |
| src_000091 | SIMD代数版本总结.md | d->zeta->轨道->slot总图 |
| src_000092 | 变换与NTT-Cooley-Tukey.md | NTT rough_note |
| src_000093 | 总结NTT作用.md | NTT本质 rough_note |
| src_000094 | 随便插值与NTT分圆上点.md | 点值表示 rough_note |

## 产出

| 类别 | 数量 |
|------|------|
| 知识卡片 | 6 (K050~K055) |
| 原子问答 | 14 (Q146~Q159) |

## FFT/DFT/NTT 质量评估

### 稳定知识
- DFT定义 (复数域，omega=e^(2pi i/n))：标准教材知识
- FFT=DFT的快速算法(Cooley-Tukey)：标准算法知识
- NTT=有限域DFT (需要n整除q-1)：标准密码学知识

### Rough note 内容 (需要作者后续确认)
1. Cooley-Tukey分解方向：笔记未给出decimation-in-time vs decimation-in-frequency的完整蝶形图
2. Bit-reversal permutation：笔记提及但未展开具体置换规则
3. NTT参数选择：q=17,n=4的toy example正确，但实际FHE参数(q,m,n)的NTT可行性未讨论
4. NTT与CKKS编码的关系 (K054)：用户笔记的类比重组，非标准定义
5. Vandermonde/DFT/Canonical Embedding的统一 (K054)：跨context综合，需核实
6. 复杂度口径：NTT的O(n log n)是普通多项式乘法的NTT加速->FHE中的同态NTT复杂度不同(需automorphism密钥)
7. 卷积定理"时域<->频域"：信号处理的直观类比正确，但FHE中不使用"频域"语义

### 不适合进入正式知识库的内容
- K053 KP5（"NTT != CKKS canonical embedding"）：断言方向正确但缺乏精确定义支撑
- K054 KP4（BGV/BFV U_ell与CKKS U_n的类比）：综合推断，非原文结论
- 以上均已标 low confidence + needs_human_review

### 只能作为 intuition 的内容
- K054 全卡：跨context综合，属于"结构启发"
- K053 KP5-KP6：NTT与CKKS的区分是基于理解和类比的概括

### 需要作者后续亲自核对
1. Cooley-Tukey在CKKS bootstrapping C2S/S2C中的确切角色
2. NTT与前几轮的U_n矩阵分解是否为同一数学框架的不同实例
3. Vandermonde at roots-of-unity的统一视角是否适用于当前论文的C2S/S2C分解
4. K053-K054 标记 low confidence 的所有声明

## 质量检查
- [x] 术语规则遵守 (C2S=coeff->slot, S2C=slot->coeff)
- [x] slot类型区分 (BGV/BFV CRT vs CKKS complex)
- [x] 变换类型区分 (DFT/FFT/NTT/canonical embedding)
- [x] NTT/FFT笔记标 rough_note + low confidence
- [x] 不把NTT等同CKKS embedding
- [x] 跨context卡标合成推断 + needs_human_review
