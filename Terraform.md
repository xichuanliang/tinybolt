# Terraform

## terraform的基础架构

- Terraform配置语言（Terraform Configuration Language）：这是Terraform的主要用户界面，使用HCL（HashiCorp Configuration Language）编写，用于描述您的基础设施资源和其关联关系。配置文件指定了您想要创建的云资源、服务、网络、存储等等。
- Terraform Core：Terraform的核心引擎，负责解析配置文件、创建资源图和执行计划。它是整个Terraform体系结构的基础，并提供了许多核心功能，如资源的依赖关系管理、状态管理和插件系统。
- 插件（Providers）：Terraform使用插件来与各种不同的云服务提供商或其他基础设施系统进行交互。每个云提供商（例如AWS、Azure、GCP等）都有自己的插件，用于管理和操作其提供的资源。
- 状态文件（State File）：Terraform使用一个状态文件来跟踪当前资源的状态和配置。这个文件记录了实际资源在云供应商中的状态，以便Terraform可以比较当前状态与配置文件中描述的目标状态之间的差异，并决定执行何种操作来达到目标状态。
- 执行计划（Execution Plan）：Terraform Core根据配置文件和状态文件生成一个执行计划。执行计划是一个详细的操作计划，描述了Terraform将如何修改资源以使其与配置文件中的期望状态匹配。
- 应用和销毁（Apply and Destroy）：当执行计划生成后，可以将其应用到基础设施中，以创建、更新或删除资源，使其达到期望的配置状态。应用操作会更新状态文件。相反，销毁操作会销毁所有资源，并将其从状态文件中移除。
- 资源依赖关系图（Resource Dependency Graph）：在执行计划过程中，Terraform会根据配置文件中的资源依赖关系构建一个资源图，以确定资源之间的关联关系和正确的创建或销毁顺序。


internal\command\apply.go Run op, err := c.RunOperation(be, opReq)

internal\command\meta.go RunOperation  op, err := b.Operation(context.Background(), opReq)

internal\backend\local\backend.go Operation方法   f = b.opApply

internal\backend\local\backend_apply.go  opApply    applyState, applyDiags = lr.Core.Apply(plan, lr.Config)

internal\terraform\context_apply.go apply  graph, operation, diags := c.applyGraph(plan, config, true)

internal\terraform\context_apply.go applyGraph  graph, moreDiags := (&ApplyGraphBuilder{

​    Config:       config,

​    Changes:       plan.Changes,

​    State:        plan.PriorState,

​    RootVariableValues: variables,

​    Plugins:       c.plugins,

​    Targets:       plan.TargetAddrs,

​    ForceReplace:    plan.ForceReplaceAddrs,

​    Operation:      operation,

​    ExternalReferences: plan.ExternalReferences,

  }).Build(addrs.RootModuleInstance)

internal\terraform\context_apply.go  Apply   walker, walkDiags := c.walk(graph, operation, &graphWalkOpts{

​    Config:   config,

​    InputState: workingState,

​    Changes:   plan.Changes,



​    // We need to propagate the check results from the plan phase,

​    // because that will tell us which checkable objects we're expecting

​    // to see updated results from during the apply step.

​    PlanTimeCheckResults: plan.Checks,



​    // We also want to propagate the timestamp from the plan file.

​    PlanTimeTimestamp: plan.Timestamp,

  })

internal\terraform\context_walk.go walk  diags := graph.Walk(walker)

internal\terraform\graph.go Walk return g.walk(walker)

internal\terraform\graph.go walk diags = diags.Append(walker.Execute(vertexCtx, ev))

internal\terraform\graph_walk.go GraphWalker Execute(EvalContext, GraphNodeExecutable) tfdiags.Diagnostics

internal\terraform\graph_walk_context.go Execute return n.Execute(ctx, w.Operation)

internal\terraform\execute.go GraphNodeExecutable Execute(EvalContext, walkOperation) tfdiags.Diagnostics

internal\terraform\node_resource_apply_instance.go Execute  return n.managedResourceExecute(ctx)

internal\terraform\node_resource_apply_instance.go managedResourceExecute state, applyDiags := n.apply(ctx, state, diffApply, n.Config, repeatData, n.CreateBeforeDestroy())

internal\terraform\node_resource_abstract_instance.go apply   resp := provider.ApplyResourceChange(providers.ApplyResourceChangeRequest{

​    TypeName:    n.Addr.Resource.Resource.Type,

​    PriorState:   unmarkedBefore,

​    Config:     unmarkedConfigVal,

​    PlannedState:  unmarkedAfter,

​    PlannedPrivate: change.Private,

​    ProviderMeta:  metaConfigVal,

  })

internal\plugin\grpc_provider.go ApplyResourceChange 