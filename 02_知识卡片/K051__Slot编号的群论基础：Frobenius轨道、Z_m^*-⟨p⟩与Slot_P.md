# Slot编号的群论基础：Frobenius轨道、Z_m^*/⟨p⟩与Slot Permutation

**card_id**: K051 | **batch**: Batch-05a3 | **confidence**: medium
**source_id**: src_000081, src_000082, src_000086, src_000091 | **claim_source_type**: 用户笔记

## Knowledge Points

1. d=ord_m(p)：最小正整数使p^d≡1(mod m)。在F_{p^d}中存在m次本原单位根ζ，Frobenius自同构Fr_p:x↦x^p作用在ζ^h上产生轨道ζ^h,ζ^{hp},...,ζ^{hp^{d-1}}。轨道长度=d。
2. 每条Frobenius轨道↔一个不可约因子F_i(X)。轨道中的根集=该因子在分裂域中的所有根。因此所有不可约因子次数相同(均为d)，因为对Φ_m的所有根，Frobenius轨道长度一致。
3. 两个指数h,h'属于同一slot当且仅当h'≡hp^k(mod m)，即属于同一⟨p⟩-陪集。Slot编号为Z_m^*/⟨p⟩的代表元集合S。槽数L=φ(m)/d。
4. 不同槽=不同的Frobenius等价类。同一槽内部包含同一等价类中的所有本原m次根(整条轨道)，不是单个根。不能混淆"槽对应一个根"和"槽对应一条轨道"。
5. Slot Permutation机制：自同构τ_j:a(X)↦a(X^j)作用在根标签上：h→hj。在slot层面诱导[h]=h⟨p⟩→[hj]=hj⟨p⟩。这是合法的slot置换，因为(hp^k)j=hj·p^k，陪集独立于代表元选择。
6. 多维槽索引：将S分解为S={g₁^{e₁}···g_t^{e_t}}，每个槽用坐标(e₁,...,e_t)标号。One-dimensional rotation ρ_i^v只动第i个坐标(e_i→e_i+v)。槽本身仍是同一批，只是编号方式不同。

**交叉引用**: K002(Y-X): 分圆多项式CRT分解的群论基础；K046(Y-X): M→U_ℓ→M^{-1}管线中的slot坐标变换；K047(Y-X): M矩阵的槽内基统一化（同一slot内部的工作）。K051从群论解释slot编号的由来。

**备注**: BGV/BFV context (Z_m^*/⟨p⟩, CRT slot)。CKKS的slot编号来自Z_m^*/⟨5⟩或类似结构（典范嵌入坐标），代数结构类似但上下文不同。
