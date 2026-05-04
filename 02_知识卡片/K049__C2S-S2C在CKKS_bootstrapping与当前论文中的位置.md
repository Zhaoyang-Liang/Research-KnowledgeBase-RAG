# C2S/S2C在CKKS bootstrapping与当前论文中的位置

**card_id**: K049 | **batch**: Batch-05a2
**source_id**: synthesized_from_user_notes | **claim_source_type**: 综合推导（基于用户笔记+PaCo论文+Bootstrapping讲义）

## Knowledge Points

1. CKKS bootstrapping管线：ModRaise→C2S→EvalMod→S2C→ModSwitch。C2S将密文的slot表示转为coefficient表示（供EvalMod在系数侧求值）；S2C将结果转回slot表示。
2. C2S/S2C≠CKKS encoding/decoding。Encoding是外部消息→多项式（σ^{-1}），C2S/S2C是同一已加密环元素内部的表示转换。
3. 当前论文主线：Split_k→(C2S_n→EvalMod_n→S2C_n)×k→Merge_k。C2S/S2C是每个leaf bootstrapping core的入口/出口步骤。
4. PaCo的加速核心：将C2S/S2C的U_n矩阵通过Cooley-Tukey分解为logN层稀疏矩阵，复杂度从O(N²)→O(NlogN)。这与当前论文的C2S_n/S2C_n直接相关。
5. C2S/S2C的线性性质是论文能使用Matrix Decomposition/Vandermonde方法的前提。Split_k不破坏该线性结构（每个leaf独立）。
6. 当前论文注意事项：(1)C2S/S2C的正确性依赖于Split/Merge的compatibility；(2)notation需与PaCo的U_n矩阵定义对齐；(3)使用coefficient split→每个leaf的C2S维度降为原来的1/k。

**交叉引用**: K039(Y-X): CKKS编码管线(σ^{-1})；K044(Y-X): C2S/S2C通用定义。K049定位C2S/S2C在bootstrapping+论文中的位置。
