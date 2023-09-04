### AsyncProcessor

A Processor to process multiple sets of input data via asyncio module.
Common pattern:
```
    processor = AsyncProcessor(PrintProcessor())
    inputData = ['doc0', 'doc1', 'doc2']
    data = self.processor.process(inputData)
```

Mostly used in [SAXSWAXS](https://github.com/CHESSComputing/CHAPBookWorkflows/tree/main/saxswaxs) workflow.
