package main

import (
	"fmt"

	"github.com/cirius-go/codegen"
	"github.com/cirius-go/generic/slice"
)

func NewAPIModule() *codegen.Seq {
	seq := codegen.NewSeq("NewAPIModule")
	seq.
		AddElems(codegen.SeqElems{
			{
				Group:           "Model",
				DefName:         "InitModelFile",
				ForwardArgsFunc: forwardAllArgs,
			},
		}).
		AddElems(codegen.SeqElems{
			{
				Group:           "Repo",
				DefName:         "InitRepoFile",
				ForwardArgsFunc: forwardAllArgs,
			},
		}).
		AddElems(codegen.SeqElems{
			{
				Group:           "API",
				DefName:         "InitAPIServiceInterfaceFile",
				ForwardArgsFunc: forwardAllArgs,
			},
			{
				Group:           "API",
				DefName:         "InitAPIFile",
				ForwardArgsFunc: forwardAllArgs,
			},
		}).
		AddElems(codegen.SeqElems{
			{
				Group:           "DTO",
				DefName:         "InitDTOFile",
				ForwardArgsFunc: forwardAllArgs,
			},
		}).
		AddElems(codegen.SeqElems{
			{
				Group:   "Service",
				DefName: "InitServiceFile",
			},
		})

	return seq
}

func NewSelectAPIHandler() *codegen.Seq {
	seq := codegen.NewSeq("NewAPIHandler")
	seq.
		AddElems(codegen.SeqElems{
			{
				Group:           "DTO",
				DefName:         "InitDTOFile",
				ForwardArgsFunc: forwardAllArgs,
			},
			{
				Group:           "DTO",
				DefName:         "InitRqRp",
				ForwardArgsFunc: forwardAllArgs,
			},
		}).
		Select(
			slice.Map(func(e string) *codegen.Seq {
				return codegen.NewSeq(e).AddElems(codegen.SeqElems{
					{
						Group:           "API",
						DefName:         fmt.Sprintf("InitAPIHandler_%s", e),
						ForwardArgsFunc: forwardAllArgs,
					},
				})
			}, "GET", "LIST", "POST", "PATCH", "DELETE"),
		).
		AddElems(codegen.SeqElems{
			{
				Group:   "Service",
				DefName: "InitServiceHandler",
			},
		})

	return seq
}

func NewInternalModule() *codegen.Seq {
	seq := codegen.NewSeq("NewInternalModule")
	seq.
		AddElems(codegen.SeqElems{
			{
				Group:           "DTO",
				DefName:         "InitDTOFile",
				ForwardArgsFunc: forwardAllArgs,
			},
		}).
		AddElems(codegen.SeqElems{
			{
				Group:   "Service",
				DefName: "InitServiceFile",
			},
		})

	return seq
}

func NewAPIInternalModule() *codegen.Seq {
	seq := codegen.NewSeq("NewAPIInternalModule")
	seq.
		AddElems(codegen.SeqElems{
			{
				Group:           "DTOImpl",
				DefName:         "InitRqRp",
				ForwardArgsFunc: forwardAllArgs,
			},
		}).
		AddElems(codegen.SeqElems{
			{
				Group:   "ServiceImpl",
				DefName: "InitServiceHandler",
			},
		})

	return seq
}
